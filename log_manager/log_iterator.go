package log_manager

import fm "file_manager"

/*
LogIterator用于遍历给定区块内的记录,由于记录从底部往上写，因此记录1,2,3,4写入后在区块的排列为
4,3,2,1，因此LogIterator会从上往下遍历记录，于是得到的记录就是4,3,2,1
*/

type LogIterator struct {
	file_manager *fm.FileManager
	blk          *fm.BlockId
	p            *fm.Page
	current_pos  uint64
	boundary     uint64
}

func NewLogIterator(file_manager *fm.FileManager, blk *fm.BlockId) *LogIterator {
	it := LogIterator{
		file_manager: file_manager,
		blk:          blk,
	}

	it.p = fm.NewPageBySize(file_manager.BlockSize())
	err := it.moveToBlock(blk)
	if err != nil {
		return nil
	}
	return &it
}

func (l *LogIterator) moveToBlock(blk *fm.BlockId) error {
	_, err := l.file_manager.Read(blk, l.p)
	if err != nil {
		return err
	}

	l.boundary = l.p.GetInt(0)
	l.current_pos = l.boundary
	return nil
}

func (l *LogIterator) Next() []byte {
	//先读取最新日志，也就是编号大的，然后依次读取编号小的
	if l.current_pos == l.file_manager.BlockSize() {
		l.blk = fm.NewBlockId(l.blk.FileName(), l.blk.Number()-1)
		err := l.moveToBlock(l.blk)
		if err != nil {
			return nil
		}
	}

	record := l.p.GetBytes(l.current_pos)
	l.current_pos += UINT64_LEN + uint64(len(record))

	return record
}

func (l *LogIterator) HasNext() bool {
	return l.current_pos < l.file_manager.BlockSize() || l.blk.Number() > 0
}
