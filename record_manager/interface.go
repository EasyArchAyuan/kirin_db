package record_manager

import fm "file_manager"

// SchemaInterface 元数据接口
type SchemaInterface interface {
	AddField(field_name string, field_type FIELD_TYPE, length int) //添加字段
	AddIntField(field_name string)                                 //添加Int类型字段
	AddStringField(field_name string, length int)                  //添加string类型字段
	Add(field_name string, sch SchemaInterface)                    //加入schema到表中
	AddAll(sch SchemaInterface)                                    //添加所有schema
	Fields() []string                                              //字段列表
	HasFields(field_name string) bool                              //是否包含某字段
	Type(field_name string) FIELD_TYPE                             //字段类型
	Length(field_name string) int                                  //字段长度
}

// LayoutInterface 字段信息接口
type LayoutInterface interface {
	Schema() SchemaInterface      //Schema对象
	Offset(field_name string) int //偏移量
	SlotSize() int                //字段长度
}

// RecordManagerInterface 记录管理器接口
type RecordManagerInterface interface {
	Block() *fm.BlockId                                //返回记录所在页面对应的区块
	GetInt(slot int, field_name string) int            //根据给定字段名取出其对应的int值
	SetInt(slot int, field_name string, val int)       //设定指定字段名的int值
	GetString(slot int, field_name string) string      //根据给定字段名获取其字符串内容
	SetString(slot int, field_name string, val string) //设置给定字段名的字符串内容
	Format()                                           //将所有插槽中的记录设定为默认值
	Delete(slot int)                                   //将给定插槽的占用标志位设置为0
	NextAfter(slot int)                                //查找给定插槽之后第一个占用标志位为1的记录
	InsertAfter(slot int)                              //查找给定插槽之后第一个占用标志位为0的记录
}

// RIDInterface 一条记录包含的区块号和插槽号
type RIDInterface interface {
	BlockNumber() int //记录所在的区块号
	Slot() int        //记录的插槽号
}

// TableScanInterface 遍历给定表的记录接口
type TableScanInterface interface {
	Close()
	HasField(field_name string) bool
	BeforeFirst()             //将指针放在第一条记录前
	Next() bool               //读取下一条记录
	MoveToRid(r RIDInterface) //跳转到给定目录
	Insert()

	GetInt(field_name string) int
	GetString(field_name string) string
	SetInt(field_name string, val int)
	SetString(field_name string, val string)
	CurrentRID() RIDInterface
	Delete()
}
