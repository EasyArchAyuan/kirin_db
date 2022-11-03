package tx

import (
	fm "file_manager"
	lg "log_manager"
	"math"
)

// CheckPointRecord <CHECKPONT>
type CheckPointRecord struct {
}

func NewCheckPointRecord() *CheckPointRecord {
	return &CheckPointRecord{}
}

func (c *CheckPointRecord) Op() RECORD_TYPE {
	return CHECKPOINT
}

func (c *CheckPointRecord) TxNumber() uint64 {
	return math.MaxUint64 //它没有对应的事务ID
}

func (c *CheckPointRecord) Undo() {
}

func (c *CheckPointRecord) ToString() string {
	return "<CHECKPOINT>"
}

func WriteCheckPointToLog(lgmr *lg.LogManager) (uint64, error) {
	rec := make([]byte, UINT64_LENGTH)
	p := fm.NewPageByBytes(rec)
	p.SetInt(0, uint64(CHECKPOINT))
	return lgmr.Append(rec)
}
