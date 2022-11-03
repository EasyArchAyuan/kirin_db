package tx

import (
	fm "file_manager"
	"fmt"
	lg "log_manager"
)

// SetIntRecord 写入整形数值：<SETINT, 0, junk, 33, 8, 542, 543>
type SetIntRecord struct {
	tx_num uint64      //事务ID
	offset uint64      //偏移量
	val    uint64      //写入的值
	blk    *fm.BlockId //写入区块ID
}

func NewSetIntRecord(p *fm.Page) *SetIntRecord {
	tpos := uint64(UINT64_LENGTH)
	tx_num := p.GetInt(tpos)
	fpos := tpos + UINT64_LENGTH
	filename := p.GetString(fpos)
	bpos := fpos + p.MaxLengthForString(filename)
	blknum := p.GetInt(bpos)
	blk := fm.NewBlockId(filename, blknum)
	opos := bpos + UINT64_LENGTH
	offset := p.GetInt(opos)
	vpos := opos + UINT64_LENGTH
	val := p.GetInt(vpos) //将日志中的字符串再次写入给定位置

	return &SetIntRecord{
		tx_num: p.GetInt(tx_num),
		offset: p.GetInt(offset),
		val:    val,
		blk:    blk,
	}
}

func (s *SetIntRecord) Op() RECORD_TYPE {
	return SETINT
}

func (s *SetIntRecord) TxNumber() uint64 {
	return s.tx_num
}

func (s *SetIntRecord) ToString() string {
	return fmt.Sprintf("<SETINT %d %d %d %s>", s.tx_num, s.blk.Number(), s.offset, s.val)
}

func (s *SetIntRecord) Undo(tx TransactionInterface) {
	tx.Pin(s.blk)
	tx.SetInt(s.blk, s.offset, s.val, false) //将原来的字符串写回去
	tx.UnPin(s.blk)
}

func WriteSetIntLog(log_manager *lg.LogManager, tx_num uint64, blk *fm.BlockId, offset uint64, val uint64) (uint64, error) {
	tpos := uint64(UINT64_LENGTH)
	fpos := tpos + UINT64_LENGTH
	p := fm.NewPageBySize(1)
	bpos := fpos + p.MaxLengthForString(blk.FileName())
	opos := bpos + UINT64_LENGTH
	vpos := opos + UINT64_LENGTH
	rec_len := vpos + UINT64_LENGTH
	rec := make([]byte, rec_len)

	p = fm.NewPageByBytes(rec)
	p.SetInt(0, uint64(SETSTRING))
	p.SetInt(tpos, tx_num)
	p.SetString(fpos, blk.FileName())
	p.SetInt(bpos, blk.Number())
	p.SetInt(opos, offset)
	p.SetInt(vpos, val)

	return log_manager.Append(rec)
}
