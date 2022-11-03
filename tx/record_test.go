package tx

import (
	fm "file_manager"
	"fmt"
	"github.com/stretchr/testify/require"
	lm "log_manager"
	"testing"
)

func TestStartRecord(t *testing.T) {
	file_manager, _ := fm.NewFileManager("recordtest", 400)
	log_manager, _ := lm.NewLogManager(file_manager, "record_file")

	tx_num := uint64(13) //事务ID
	p := fm.NewPageBySize(32)
	p.SetInt(0, uint64(START))
	p.SetInt(8, tx_num)
	start_record := NewStartRecord(p, log_manager)
	expected_str := fmt.Sprintf("<START %d>", tx_num)
	require.Equal(t, expected_str, start_record.ToString())

	_, err := start_record.WriteToLog()
	require.Nil(t, err)

	//iter := log_manager.Iterator()
	////检查写入的日志是否符号预期
	//rec := iter.Next()
	//rec_op := binary.LittleEndian.Uint64(rec[0:8])
	//rec_tx_num := binary.LittleEndian.Uint64(rec[8:len(rec)])
	//require.Equal(t, rec_op, START)
	//require.Equal(t, rec_tx_num, tx_num)
}
