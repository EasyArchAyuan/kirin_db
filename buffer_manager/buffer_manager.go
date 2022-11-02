package buffer_manager

import (
	"errors"
	fm "file_manager"
	lm "log_manager"
	"sync"
	"time"
)

const (
	MAX_TIME = 3 //分配页面时最多等待3秒
)

type BufferManager struct {
	buffer_pool   []*Buffer  //多个buffer组成内存池
	num_available uint32     //可用buffer数量
	mu            sync.Mutex //互斥锁
}

func NewBufferManager(fm *fm.FileManager, lm *lm.LogManager, num_buffers uint32) *BufferManager {
	buffer_manager := &BufferManager{num_available: num_buffers}
	for i := uint32(0); i < num_buffers; i++ {
		buffer := NewBuffer(fm, lm)
		buffer_manager.buffer_pool = append(buffer_manager.buffer_pool, buffer)
	}
	return buffer_manager
}

// Available 当前可用缓存页面数量
func (b *BufferManager) Available() uint32 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.num_available
}

func (b *BufferManager) FlushAll(txnum int32) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, buff := range b.buffer_pool {
		if buff.ModifyingTx() == txnum {
			buff.Flush()
		}
	}
}

func (b *BufferManager) Pin(blk *fm.BlockId) (*Buffer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	start := time.Now()
	buffer := b.tryPin(blk)
	//如果无法获得缓存页面，那么让调用者等待一段时间后再次尝试
	for buffer != nil && !b.waitingTooLong(start) {
		time.Sleep(MAX_TIME * time.Second)
		buffer = b.tryPin(blk)
		if buffer == nil {
			return nil, errors.New("no buffer available,careful for dead lock")
		}
	}
	return buffer, nil
}

func (b *BufferManager) Unpin(buff *Buffer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if buff == nil {
		return
	}

	buff.UnPin()
	if !buff.IsPinned() {
		b.num_available += 1
		//todo: 唤醒所有等待它的线程,notifyAll()
	}
}

// waitingTooLong 超时机制,解决死锁,超过3s判断
func (b *BufferManager) waitingTooLong(start time.Time) bool {
	seconds := time.Since(start).Seconds()
	if seconds >= MAX_TIME {
		return true
	}
	return false
}

// tryPin 看给定的区块是否已经被读入某个缓存页,查看是否还有可用缓存页，然后将区块数据写入
func (b *BufferManager) tryPin(blk *fm.BlockId) *Buffer {
	buffer := b.findExistingBuffer(blk)
	if buffer == nil {
		buffer = b.chooseUnpinBuffer()
		if buffer == nil {
			return nil
		}
		buffer.AssignToBlock(blk)
	}
	//可用buffer数量 -1
	if !buffer.IsPinned() {
		b.num_available -= 1
	}
	buffer.Pin()
	return buffer
}

// findExistingBuffer 查看当前请求的区块是否已经被加载到了某个缓存页，如果是，那么直接返回即可
func (b *BufferManager) findExistingBuffer(blk *fm.BlockId) *Buffer {
	for _, buffer := range b.buffer_pool {
		block := buffer.Block()
		if block != nil && block.Equals(blk) {
			return buffer
		}
	}
	return nil
}

// chooseUnpinBuffer 选取一个没有被使用的缓存页
func (b *BufferManager) chooseUnpinBuffer() *Buffer {
	for _, buffer := range b.buffer_pool {
		if !buffer.IsPinned() {
			return buffer
		}
	}
	return nil
}
