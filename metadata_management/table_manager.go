package metadata_management

import (
	rm "record_manager"
	"tx"
)

const MAX_NAME = 16

// TableManager metadata表
type TableManager struct {
	tcatLayout *rm.Layout
	fcatLayout *rm.Layout
}

func (t TableManager) CreateTable(tblName string, sch *rm.Schema, tx *tx.Transaction) {
	//TODO implement me
	panic("implement me")
}

func (t TableManager) GetLayout(tblName string, tx *tx.Transaction) *rm.Layout {

	//TODO implement me
	panic("implement me")
}

func NewTableManager(isNew bool, tx *tx.Transaction) *TableManager {
	//创建两个表专门用于存储新建数据库表的元数据
	tableMgr := &TableManager{}

	tcatSchema := rm.NewSchema()                   //存储表的元数据
	tcatSchema.AddStringField("tblname", MAX_NAME) //表的名称
	tcatSchema.AddIntField("slotsize")             //表一条记录的长度
	tableMgr.tcatLayout = rm.NewLayoutWithSchema(tcatSchema)

	fcatSchema := rm.NewSchema()                   //存储所创建表每个字段的元数据
	fcatSchema.AddStringField("tblname", MAX_NAME) //字段所在的表名字
	fcatSchema.AddStringField("tblname", MAX_NAME) //字段的名称
	fcatSchema.AddIntField("type")                 //字段类型
	fcatSchema.AddIntField("length")               //字段数据长度
	fcatSchema.AddIntField("offset")               //字段在记录中的偏移
	tableMgr.fcatLayout = rm.NewLayoutWithSchema(fcatSchema)

	//如果当前数据表是第一次创建，那么为这个表创建两个元数据表
	if isNew {
		tableMgr.CreateTable("tblcat", tcatSchema, tx)
		tableMgr.CreateTable("fldcat", fcatSchema, tx)
	}
	return tableMgr
}
