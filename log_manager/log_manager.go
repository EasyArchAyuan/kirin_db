package log_manager

import (
	fm "file_manager"
	"sync"
)

const (
	UINT64_LEN = 8
)

type LogManager struct {
	file_manager   *fm.FileManager //操作文件的管理器
	log_file       string          //日志文件名称
	log_page       *fm.Page        //存储日志的缓冲区
	current_blk    *fm.BlockId     //日志当前写入的区块号
	latest_lsn     uint64          //当前日志序列号
	last_saved_lsn uint64          //上次存储到磁盘的日志序列号
	mu             sync.Mutex      //互斥锁
}

// 日志末尾追加写
func (l *LogManager) appendNewBlock() (*fm.BlockId, error) {
	//当缓冲器用完后调用该接口分配新内存
	blk, err := l.file_manager.Append(l.log_file) //在日志二进制文件末尾添加一个区块
	if err != nil {
		return nil, err
	}
	/*
		添加日志时从内存的底部往上走，例如内存400字节，日志100字节，那么
		日志将存储在内存的300到400字节处，因此我们需要把当前内存可用底部偏移
		写入头8个字节
	*/
	l.log_page.SetInt(0, uint64(l.file_manager.BlockSize()))
	_, err = l.file_manager.Write(&blk, l.log_page)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}

func NewLogManager(file_manager *fm.FileManager, log_file string) (*LogManager, error) {
	log_mgr := LogManager{
		file_manager:   file_manager,
		log_file:       log_file,
		log_page:       fm.NewPageBySize(file_manager.BlockSize()),
		last_saved_lsn: 0,
		latest_lsn:     0,
	}

	log_size, err := file_manager.Size(log_file)
	if err != nil {
		return nil, err
	}

	if log_size == 0 {
		blk, err := log_mgr.appendNewBlock()
		if err != nil {
			return nil, err
		}
		log_mgr.current_blk = blk
	} else {
		log_mgr.current_blk = fm.NewBlockId(log_mgr.log_file, log_size-1)
		_, err := file_manager.Read(log_mgr.current_blk, log_mgr.log_page)
		if err != nil {
			return nil, err
		}
	}

	return &log_mgr, nil
}

func (l *LogManager) Flush() error {
	//将当前日志写入磁盘
	_, err := l.file_manager.Write(l.current_blk, l.log_page)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogManager) FlushByLSN(lsn uint64) error {
	if lsn > l.last_saved_lsn {
		//将当前日志写入磁盘
		err := l.Flush()
		if err != nil {
			return err
		}
		l.last_saved_lsn = lsn
	}

	return nil
}

func (l *LogManager) Append(log_record []byte) (uint64, error) {
	//添加日志
	l.mu.Lock()
	defer l.mu.Unlock()

	//获得可写入的底部偏移
	boundary := l.log_page.GetInt(0)
	record_size := uint64(len(log_record))
	bytes_need := record_size + UINT64_LEN
	var err error
	if int(boundary-bytes_need) < UINT64_LEN {
		//当前容量不够,现将当前日志写入磁盘
		err = l.Flush()
		if err != nil {
			return l.latest_lsn, err
		}
		//生成新区块用于写新数据
		l.current_blk, err = l.appendNewBlock()
		if err != nil {
			return l.latest_lsn, err
		}

		boundary = l.log_page.GetInt(0)
	}

	record_pos := boundary - bytes_need         //我们从底部往上写入
	l.log_page.SetBytes(record_pos, log_record) //设置下次可以写入的位置
	l.log_page.SetInt(0, record_pos)
	l.latest_lsn += 1 //记录新加入日志的编号

	return l.latest_lsn, err
}

func (l *LogManager) Iterator() *LogIterator {
	l.Flush()
	return NewLogIterator(l.file_manager, l.current_blk)
}
