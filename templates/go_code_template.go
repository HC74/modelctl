package templates

// 仓储
var warehouseCode = ` // Package template Generation date {{ now }}
package {{.package_name}}
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var DB *gorm.DB
func GormMysql() *gorm.DB {
	mysqlConfig := mysql.Config{
		DSN:                       "{{ .dns }}",    // DSN data source name
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig))
	if err != nil {
		// TODO
	}
	DB = db
	return DB
}
`

// GetWarehouseCode 获取仓储实体
func GetWarehouseCode() string {
	return warehouseCode
}

func GetGeneratorStructCode() string {
	return MODEL_TEMPLATE
}

const (
	MODEL_TEMPLATE = `{{- $packageName := .Package}}{{$structName := .StructName}}{{$tableName := .TableName}}{{$hasTime := .HasTime -}}
{{"// Package "}}{{$packageName}}{{" This file is generated by Cli, please do not modify"}}
package {{$packageName}}
{{if ($hasTime)}}
import "time"
{{end }}
type {{$structName}} struct {
{{- range $i, $v := .Cols}}
	{{$v.Column}}	{{$v.ColumnType}}	{{$v.Tag }}		{{$v.Remark -}}
{{end}}
}

// TableName 表名
func (*{{$structName}}) TableName() string {
    return "{{$tableName}}"
}`
)
