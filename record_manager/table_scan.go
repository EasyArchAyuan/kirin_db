package record_manager

import "tx"

// TableScan 表扫描类（不直接处理tx，blockId等底层数据）
type TableScan struct {
	tx           *tx.Transaction        //事务接口
	layout       LayoutInterface        //视图接口
	rp           RecordManagerInterface //记录管理器接口
	file_name    string                 //文件名称
	current_slot int                    //当前区块
}

func NewTableScan(tx *tx.Transaction, table_name string, layout LayoutInterface) *TableScan {
	table_scan := &TableScan{
		tx:        tx,
		layout:    layout,
		file_name: table_name + ".tbl",
	}

	size, err := tx.Size(table_scan.file_name)
	if err != nil {
		panic(err)
	}

	if size == 0 {
		//如果文件为空则增加一个区块
		table_scan.MoveToNewBlock()
	} else {
		//先读取第一个区块
		table_scan.MoveToBlock(0)
	}
	return table_scan
}

func (t TableScan) Close() {
	if t.rp != nil {
		t.tx.UnPin(t.rp.Block())
	}
}

func (t TableScan) HasField(field_name string) bool {
	return t.layout.Schema().HasFields(field_name)
}

func (t TableScan) BeforeFirst() {
	t.MoveToBlock(0)
}

func (t TableScan) Next() bool {
	//如果在当前区块找不到给定有效记录则遍历后续区块，直到所有区块都遍历为止
	t.current_slot = t.rp.NextAfter(t.current_slot)
	for t.current_slot < 0 {
		if t.AtLastBlock() {
			//直到最后一个区块都找不到给定插槽
			return false
		}
		t.MoveToBlock(int(t.rp.Block().Number() + 1))
		t.current_slot = t.rp.NextAfter(t.current_slot)
	}
	return true
}

func (t TableScan) MoveToRid(r RIDInterface) {
	//TODO implement me
	panic("implement me")
}

func (t TableScan) Insert() {
	//TODO implement me
	panic("implement me")
}

func (t TableScan) GetInt(field_name string) int {
	return t.rp.GetInt(t.current_slot, field_name)
}

func (t TableScan) GetString(field_name string) string {
	return t.rp.GetString(t.current_slot, field_name)
}

func (t TableScan) SetInt(field_name string, val int) {
	t.rp.SetInt(t.current_slot, field_name, val)

}

func (t TableScan) SetString(field_name string, val string) {
	t.rp.SetString(t.current_slot, field_name, val)
}

func (t TableScan) CurrentRID() RIDInterface {
	//TODO implement me
	panic("implement me")
}

func (t TableScan) Delete() {
	//TODO implement me
	panic("implement me")
}

func (t TableScan) MoveToNewBlock() {

}

func (t TableScan) MoveToBlock(i int) {

}
