package model

import (
	"github.com/HC74/modelctl/utils"
	"strings"
)

var (
	// go类型映射关系
	goMapper   = make(map[string]string, 55)
	msGoMapper = make(map[string]string, 10)
)

func init() {
	// 初始化映射关系
	//布尔
	goMapper["boolean"] = "bool"
	// 小int
	goMapper["tinyint"] = "int8"

	goMapper["smallint"] = "int16"
	goMapper["year"] = "int16"
	//整数类型
	goMapper["integer"] = "int32"
	goMapper["mediumint"] = "int32"
	goMapper["int"] = "int32"
	// long
	goMapper["bigint"] = "int64"

	// 日期
	goMapper["date"] = "time.Time"
	goMapper["timestamp without time zone"] = "time.Time"
	goMapper["timestamp with time zone"] = "time.Time"
	goMapper["time with time zone"] = "time.Time"
	goMapper["time without time zone"] = "time.Time"
	goMapper["timestamp"] = "time.Time"
	goMapper["datetime"] = "time.Time"
	goMapper["time"] = "time.Time"

	// 字节
	goMapper["bytea"] = "[]byte"
	goMapper["binary"] = "[]byte"
	goMapper["varbinary"] = "[]byte"
	goMapper["tinyblob"] = "[]byte"
	goMapper["blob"] = "[]byte"
	goMapper["mediumblob"] = "[]byte"
	goMapper["longblob"] = "[]byte"

	//字符串
	goMapper["text"] = "string"
	goMapper["character"] = "string"
	goMapper["character varying"] = "string"
	goMapper["tsvector"] = "string"
	goMapper["bit"] = "string"
	goMapper["bit varying"] = "string"
	goMapper["money"] = "string"
	goMapper["json"] = "string"
	goMapper["jsonb"] = "string"
	goMapper["xml"] = "string"
	goMapper["point"] = "string"
	goMapper["interval"] = "string"
	goMapper["line"] = "string"
	goMapper["ARRAY"] = "string"
	goMapper["char"] = "string"
	goMapper["varchar"] = "string"
	goMapper["tinytext"] = "string"
	goMapper["mediumtext"] = "string"
	goMapper["longtext"] = "string"

	// 单浮点数
	goMapper["real"] = "float32"

	// 双浮点
	goMapper["numeric"] = "float64"
	goMapper["decimal"] = "float64"
	goMapper["double precision"] = "float64"
	goMapper["float"] = "float64"
	goMapper["double"] = "float64"

	// mssql
	msGoMapper["int"] = "int"
	msGoMapper["bigint"] = "int64"
	msGoMapper["datetime2"] = "time.Time"
	msGoMapper["datetime"] = "time.Time"
	msGoMapper["smalldatetime"] = "time.Time"
	msGoMapper["date"] = "time.Time"
	msGoMapper["bit"] = "bool"
	msGoMapper["text"] = "string"

}

// ColumnMetaData 列元数据
type ColumnMetaData struct {
	// 列名
	Name string
	// 列类型
	DBType string
	// 对应 go语言的什么类型
	GoType string
	// 无符号？？？
	IsUnsigned bool
	// 是否空
	IsNullable bool
	// 表名称
	TableName string
	// 备注
	Remark string
	// 是否为主键
	IsPKey bool
}

func (c *ColumnMetaData) InitGoType(databaseType string) {
	c.GoType = GetGoType(c.DBType, databaseType)
}

// ColumnMetaDataList 列元数据集
type ColumnMetaDataList []*ColumnMetaData

// TableMetaData 表元数据
type TableMetaData struct {
	// 表名
	Name string
	// 列元数据
	Columns ColumnMetaDataList
}

// TableMataDataList 表元数据集
type TableMataDataList []*TableMetaData

// NewColumnMetaData 创建列数据源对象
func NewColumnMetaData(name string, isNullable bool, dataType string, isUnsigned bool, tableName, remark string, isKey bool, databaseType string) (columnMeta *ColumnMetaData) {
	columnMeta = &ColumnMetaData{
		Name:       name,
		IsNullable: isNullable,
		DBType:     dataType,
		IsUnsigned: isUnsigned,
		TableName:  tableName,
		Remark:     remark,
		IsPKey:     isKey,
	}
	columnMeta.InitGoType(databaseType)
	return
}

func GetGoType(k, t string) (v string) {
	if t == "mssql" {
		v = msGoMapper[k]
		if strings.HasPrefix(k, "varchar") || strings.HasPrefix(k, "nvarchar") {
			v = "string"
		}
		if strings.HasPrefix(k, "decimal") {
			v = "float64"
		}
	} else {
		v = goMapper[k]
	}
	if utils.IsEmpty(v) {
		v = "string"
	}
	return
}
