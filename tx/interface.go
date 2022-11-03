package tx

import fm "file_manager"

// TransactionInterface 事务对象的接口
type TransactionInterface interface {
	Commit()
	RollBack()
	Recover()
	Pin(blk *fm.BlockId)
	UnPin(blk *fm.BlockId)
	GetInt(blk *fm.BlockId, offset uint64) uint64
	SetInt(blk *fm.BlockId, offset uint64, val uint64, okToLog bool)
	GetString(blk *fm.BlockId, offset uint64) string
	SetString(blk *fm.BlockId, offset uint64, val string, okToLog bool)
	AvailableBuffers() uint64
	Size(filename string) uint64
	Append(filename string) uint64
	BlockSize() uint64
}

type RECORD_TYPE uint64

// undo的六种日志格式
const (
	CHECKPOINT RECORD_TYPE = iota
	START
	COMMIT
	ROLLBACK
	SETINT
	SETSTRING
)

const (
	UINT64_LENGTH = 8
)

type LogRecordInterface interface {
	Op() RECORD_TYPE              //返回记录的类别
	TxNumber() uint32             //对应事务的编号
	Undo(tx TransactionInterface) //回滚操作
	ToString() string             //获取记录的字符串内容
}
