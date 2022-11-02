package buffer_manager

import (
	fm "file_manager"
	log "log_manager"
)

type Buffer struct {
	fm       *fm.FileManager
	lm       *log.LogManager
	contents *fm.Page //存储磁盘数据的缓存页面
	blk      *fm.BlockId
	pins     uint32 //被引用次数
	txnum    int32  //交易号
	lsn      uint64 //对应日志号
}

func NewBuffer(fmgr *fm.FileManager, lmgr *log.LogManager) *Buffer {
	return &Buffer{
		fm:       fmgr,
		lm:       lmgr,
		contents: fm.NewPageBySize(fmgr.BlockSize()),
	}
}

func (b *Buffer) Contents() *fm.Page {
	return b.contents
}

func (b *Buffer) Block() *fm.BlockId {
	return b.blk
}

func (b *Buffer) Pin() {
	b.pins += 1
}

func (b *Buffer) UnPin() {
	b.pins -= 1
}

func (b *Buffer) SetModified(txnum int32, lsn uint64) {
	//如果客户修改了页面数据，必须调用该接口通知Buffer
	b.txnum = txnum
	if lsn > 0 {
		b.lsn = lsn
	}
}

func (b *Buffer) IsPinned() bool {
	return b.pins > 0
}

func (b *Buffer) ModifyingTx() int32 {
	return b.txnum
}

func (b *Buffer) Flush() {
	if b.txnum >= 0 {
		//代表当前页面的数据已经被修改过,需要写入磁盘啦
		err := b.lm.FlushByLSN(b.lsn)
		if err != nil {
			return
		}

		//先将修改操作对应的日志写入
		_, err = b.fm.Write(b.blk, b.Contents())
		if err != nil {
			return
		}

		b.txnum = -1
	}
}

// AssignToBlock 把当前页面分发给其他块
func (b *Buffer) AssignToBlock(block *fm.BlockId) {
	//当页面分发给新数据时需要判断当前页面数据是否需要写入磁盘
	b.Flush()
	b.blk = block
	//将对应数据从磁盘读取页面
	_, err := b.fm.Read(b.blk, b.Contents())
	if err != nil {
		return
	}
	b.pins = 0
}
