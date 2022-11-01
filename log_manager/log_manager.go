package log_manager

import (
	fm "file_manager"
	"sync"
)

const (
	UINT64_LEN = 8
)

type LogManager struct {
	file_manager   *fm.FileManager //操作文件的管理器
	log_file       string          //日志文件名称
	latest_lsn     uint64          //当前日志序列号
	last_saved_lsn uint64          //上次存储到磁盘的日志序列号
	mu             sync.Mutex      //互斥锁
}
