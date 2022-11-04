package tx

import (
	bm "buffer_manager"
	fm "file_manager"
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
	//recovery_mgr   *RecoveryManager
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
