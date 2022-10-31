package file_manager

import (
	"crypto/sha256"
	"fmt"
)

// BlockId 区块的抽象,对应二进制文件某个位置的一块连续内存的记录
type BlockId struct {
	file_name string //对应磁盘二进制文件
	blk_num   uint64 //二进制文件中区块的标号
}

func NewBlockId(file_name string, blk_num uint64) *BlockId {
	return &BlockId{
		file_name: file_name,
		blk_num:   blk_num,
	}
}

func (b *BlockId) FileName() string {
	return b.file_name
}

func (b *BlockId) Number() uint64 {
	return b.blk_num
}

func (b *BlockId) Equals(other *BlockId) bool {
	return b.file_name == other.FileName() && b.blk_num == other.blk_num
}

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (b *BlockId) HashCode() string {
	return asSha256(*b)
}
