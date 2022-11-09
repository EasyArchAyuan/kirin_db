package record_manager

type FIELD_TYPE int

const (
	INTEGER FIELD_TYPE = iota
	VARCHAR
)

type FieldInfo struct {
	field_type FIELD_TYPE //字段类型
	length     int        //字段长度
}

// Schema 描述给定表中一条记录包含哪些字段
type Schema struct {
	fields []string              //字段名称列表
	info   map[string]*FieldInfo //key:字段名称, value:字段类型,长度
}

func newFieldInfo(field_type FIELD_TYPE, length int) *FieldInfo {
	return &FieldInfo{
		field_type: field_type,
		length:     length,
	}
}

func NewSchema() *Schema {
	return &Schema{
		fields: make([]string, 0),
		info:   make(map[string]*FieldInfo),
	}
}

func (s *Schema) AddField(field_name string, field_type FIELD_TYPE, length int) {
	//字段名称添加到数组里
	s.fields = append(s.fields, field_name)
	//再添加字段信息到schema结构体里
	s.info[field_name] = newFieldInfo(field_type, length)
}

func (s *Schema) AddIntField(field_name string) {
	//对于整型字段而言，长度没有作用
	s.AddField(field_name, INTEGER, 0)
}

func (s *Schema) AddStringField(field_name string, length int) {
	s.AddField(field_name, VARCHAR, length)
}

func (s *Schema) Add(field_name string, sch SchemaInterface) {
	s.AddField(field_name, sch.Type(field_name), s.Length(field_name))
}

func (s *Schema) AddAll(sch SchemaInterface) {
	fields := sch.Fields()
	for _, value := range fields {
		s.Add(value, sch)
	}
}

func (s *Schema) Fields() []string {
	return s.fields
}

func (s *Schema) HasFields(field_name string) bool {
	for _, field := range s.Fields() {
		if field == field_name {
			return true
		}
	}
	return false
}

func (s *Schema) Type(field_name string) FIELD_TYPE {
	return s.info[field_name].field_type
}

func (s *Schema) Length(field_name string) int {
	return s.info[field_name].length
}
