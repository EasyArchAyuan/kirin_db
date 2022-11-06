package tx

import (
	bm "buffer_manager"
	"errors"
	fm "file_manager"
	"fmt"
	lm "log_manager"
	"sync"
)

var tx_num_mu sync.Mutex
var next_tx_num = int32(0)

func NxtTxNum() int32 {
	tx_num_mu.Lock()
	defer tx_num_mu.Unlock()

	next_tx_num += 1

	return next_tx_num
}

type Transaction struct {
	//concur_mgr     *ConCurrencyManager
	recovery_mgr   *RecoveryManager
	file_manager   *fm.FileManager
	log_manager    *lm.LogManager
	buffer_manager *bm.BufferManager
	my_buffers     *BufferList
	tx_num         int32
}

func NewTransaction(file_manager *fm.FileManager, log_manager *lm.LogManager, buffer_manager *bm.BufferManager) *Transaction {
	tx := &Transaction{
		file_manager:   file_manager,
		log_manager:    log_manager,
		buffer_manager: buffer_manager,
		my_buffers:     NewBufferList(buffer_manager),
		tx_num:         NxtTxNum(),
	}
	//创建同步管理器
	//tx.concur_mgr = NewConcurrencyManager()
	//创建恢复管理器
	//tx.recovery_mgr = NewRecoveryManager(tx, tx_num, log_manager, buffer_manager)
	return tx
}
func (t *Transaction) Commit() {
	//调用恢复管理器
	//t.concur_mgr.Release()
	t.recovery_mgr.Commit()
	r := fmt.Sprintf("transation %d committed", t.tx_num)
	fmt.Println(r)
	//释放缓存管理页面
	t.my_buffers.UnpinAll()

}
func (t *Transaction) RollBack() {
	//调用恢复管理器rollback
	t.recovery_mgr.Rollback()
	//t.concur_mgr.Release()
	r := fmt.Sprintf("transation %d roll back", t.tx_num)
	fmt.Println(r)
	//释放缓存管理页面
	t.my_buffers.UnpinAll()
}
func (t *Transaction) Recover() {
	//系统启动时会在所有交易执行前执行该函数
	t.buffer_manager.FlushAll(t.tx_num)
	//调用恢复管理器的recover接口
	t.recovery_mgr.Recover()
}
func (t *Transaction) Pin(blk *fm.BlockId) {
	t.my_buffers.Pin(blk)
}

func (t *Transaction) UnPin(blk *fm.BlockId) {
	t.my_buffers.Unpin(blk)
}

func (t *Transaction) buffer_no_exist(blk *fm.BlockId) error {
	err_s := fmt.Sprintf("No buffer found for given blk : %d with file name: %s\n",
		blk.Number(), blk.FileName())
	err := errors.New(err_s)
	return err
}

func (t *Transaction) GetInt(blk *fm.BlockId, offset uint64) (int64, error) {
	//调用同步管理器加s锁
	//err := t.concur_mgr.SLock(blk)
	//if err != nil {
	//	return -1, err
	//}

	buff := t.my_buffers.get_buffer(blk)
	if buff == nil {
		return -1, t.buffer_no_exist(blk)
	}
	//兼容page的uint64做个转换
	return int64(buff.Contents().GetInt(offset)), nil
}

func (t *Transaction) GetString(blk *fm.BlockId, offset uint64) (string, error) {
	//调用同步管理器加s锁
	//err := t.concur_mgr.SLock(blk)
	//if err != nil {
	//	return "", err
	//}

	buff := t.my_buffers.get_buffer(blk)
	if buff == nil {
		return "", t.buffer_no_exist(blk)
	}

	return buff.Contents().GetString(offset), nil
}

func (t *Transaction) SetInt(blk *fm.BlockId, offset uint64, val int64, okToLog bool) error {
	//调用同步管理器加x锁
	//err := t.concur_mgr.XLock(blk)
	//if err != nil {
	//	return err
	//}

	buff := t.my_buffers.get_buffer(blk)
	if buff == nil {
		return t.buffer_no_exist(blk)
	}

	var lsn uint64
	if okToLog {
		//调用恢复管理器的SetInt方法
		//lsn, err = t.recovery_mgr.SetInt(buff, offset, val)
		//if err != nil {
		//	return err
		//}
	}

	p := buff.Contents()
	p.SetInt(offset, uint64(val))
	buff.SetModified(t.tx_num, lsn)
	return nil
}

func (t *Transaction) SetString(blk *fm.BlockId, offset uint64, val string, okToLog bool) error {
	//使用同步管理器加x锁
	//err := t.concur_mgr.XLock(blk)
	//if err != nil {
	//	return err
	//}

	buff := t.my_buffers.get_buffer(blk)
	if buff == nil {
		return t.buffer_no_exist(blk)
	}

	var lsn uint64

	if okToLog {
		//调用恢复管理器SetString方法
		//lsn, err = t.recovery_mgr.SetString(buff, offset, val)
		//if err != nil {
		//	return err
		//}
	}

	p := buff.Contents()
	p.SetString(offset, val)
	buff.SetModified(t.tx_num, lsn)
	return nil
}

func (t *Transaction) Size(file_name string) (uint64, error) {
	//调用同步管理器加锁
	//dummy_blk := fm.NewBlockId(file_name, uint64(END_OF_FILE))
	//err := t.concur_mgr.SLock(dummy_blk)
	//if err != nil {
	//	return 0, err
	//}
	s, _ := t.file_manager.Size(file_name)
	return s, nil
}

func (t *Transaction) Append(file_name string) (*fm.BlockId, error) {
	//调用同步管理器加锁
	//dummy_blk := fm.NewBlockId(file_name, END_OF_FILE)
	//err := t.concur_mgr.XLock(dummy_blk)
	//if err != nil {
	//	return nil, err
	//}
	blk, err := t.file_manager.Append(file_name)
	if err != nil {
		return nil, err
	}

	return &blk, nil
}

func (t *Transaction) BlockSize() uint64 {
	return t.file_manager.BlockSize()
}

func (t *Transaction) AvailableBuffers() uint64 {
	return uint64(t.buffer_manager.Available())
}
