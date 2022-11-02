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
	txnum    uint64 //交易号
	lsn      uint64 //对应日志号
}

func NewBuffer(fmgr *fm.FileManager, lmgr *log.LogManager) *Buffer {
	return &Buffer{
		fm:       fmgr,
		lm:       lmgr,
		contents: fm.NewPageBySize(fmgr.BlockSize()),
	}
}
