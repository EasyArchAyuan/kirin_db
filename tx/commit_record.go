package tx

import (
	fm "file_manager"
	"fmt"
	lg "log_manager"
)

// CommitRecord <COMMIT, 1>
type CommitRecord struct {
	tx_num uint64
}

func NewCommitkRecordRecord(p *fm.Page) *CommitRecord {
	return &CommitRecord{
		tx_num: p.GetInt(UINT64_LENGTH),
	}
}

func (r *CommitRecord) Op() RECORD_TYPE {
	return COMMIT
}

func (r *CommitRecord) TxNumber() uint64 {
	return r.tx_num
}

func (r *CommitRecord) ToString() string {
	return fmt.Sprintf("<COMMIT %d>", r.tx_num)
}

func (c *CommitRecord) Undo() {
}

func WriteCommitkRecordLog(lgmr *lg.LogManager, tx_num uint64) (uint64, error) {
	rec := make([]byte, 2*UINT64_LENGTH)
	p := fm.NewPageByBytes(rec)
	p.SetInt(0, UINT64_LENGTH)
	p.SetInt(UINT64_LENGTH, tx_num)
	return lgmr.Append(rec)
}
