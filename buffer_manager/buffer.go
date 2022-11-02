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
	lsn      int32  //对应日志号
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
	b.pins = b.pins + 1
}

func (b *Buffer) UnPin() {
	b.pins = b.pins - 1
}

func (b *Buffer) SetModified(txnum int32, lsn int32) {
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
