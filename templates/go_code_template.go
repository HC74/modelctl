package templates

var warehouseCode = `

`

// GetWarehouseCode 获取仓储实体
func GetWarehouseCode() string {
	return warehouseCode
}

// GetGeneratorStructCode 获取结构体模板
func GetGeneratorStructCode() string {
	return MODEL_TEMPLATE
}

// GetCURDCode 获取curd模板
func GetCURDCode() string {
	return CRUD_TEMPLATE
}

// GetWAREHOUSECode 获取仓储模板
func GetWAREHOUSECode(t string) string {
	if t == "mssql" {
		return WAREHOUSE_SQLSERVER_TEMPLATE
	}
	return WAREHOUSE_TEMPLATE
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
	CRUD_TEMPLATE = `{{- $structName := .StructName}}{{ $dbName := .DBName }}{{ $hasKey := .HasKey -}}
{{- $idType:= .KeyType }}{{ $packageName := .Package }}{{$keyStructColumn := .KeyStructColumn}}{{$keyColumn := .KeyColumn}}
{{"// Package "}}{{$packageName}}{{" This file is automatically generated, will not be overwritten, and can be modified"}}
package {{$packageName}}

import (
    "reflect"
    "strings"
)

// InsertOne 插入一条数据
func (data *{{$structName}}) InsertOne() error {
    err := {{$dbName}}.Create(&data).Error
    if err != nil {
        return err
    }
    return nil
}

// Insert{{$structName}}Many 批量插入 @param batch:每批的数量
func Insert{{$structName}}Many(datas []*{{$structName}}, batch int) error {
    err := {{$dbName}}.CreateInBatches(datas, batch).Error
    if err != nil {
        return err
    }
    return nil
}
{{if ($hasKey)}}
// FindById 根据ID查找 @param id:唯一标识
func (data *{{$structName}}) FindById() error {
    err := {{$dbName}}.Where("{{$keyColumn}} = ?", data.{{$keyStructColumn}}).First(&data).Error
    if err != nil {
        return err
    }
    return nil
}

// UpdateById 根据ID查找 @param id:唯一标识
func (data *{{$structName}}) UpdateById() error {
    err := {{$dbName}}.Save(&data).Error
    if err != nil {
        return err
    }
    return nil
}
{{end}}


// Paging{{$structName}} 分页
// @param pageNum: 页数
// @param pageSize: 每页大小
// @param orders: 排序
// @param maps: 条件
// @param args: 值
// @returns r1: 数据集 r2: 总页数 r3: 总数据量 r4: 异常
// maps and args 例如 maps: id = ? and args: 1
func Paging{{$structName}}(pageNum, pageSize int, orders string, maps interface{}, args ...interface{}) ([]*{{$structName}}, int, int64, error) {
    var (
        results    []*{{$structName}}
        err   error
        size  = pageSize
        total int64
    )
    if pageNum == 0 {
        panic("页数不能为0")
    }
    pageNum--
    if strings.TrimSpace(orders) == "" {
        {{$dbName}} = {{$dbName}}.Order(orders)
    }
    if reflect.ValueOf(maps).IsValid() && maps != "" {
        {{$dbName}} = {{$dbName}}.Where(maps, args...)
    }
    err = {{$dbName}}.Model({{$structName}}{}).Count(&total).Error
    if err != nil {
        return nil, 0, 0, err
    }
    totalPageNum := total / int64(pageSize)
    if total%int64(pageSize) != 0 {
        totalPageNum++
    }
    err = {{$dbName}}.Offset(pageNum * size).Limit(size).Find(&results).Error
    if err != nil {
        return nil, 0, 0, err
    }
    return results, int(totalPageNum), total, nil
}
`
	WAREHOUSE_TEMPLATE = `{{- $packageName := .Package}}{{ $dbName := .DbName}}{{ $user := .User}}{{ $password := .Password -}}
{{- $host := .Host }}{{ $now := .Now }}{{ $dbType := .DbType -}}
{{"// Package "}}{{$packageName}}{{"This model file is automatically generated in "}}{{$now}}{{" and can be edited at will, and it will not be overwritten when it is generated again"}}
package {{$packageName}}
import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "time"
)

var DB_{{- $dbType }} *gorm.DB

func Setup() {
    var (
        dbName, user, password, host, dsn string
        err                               error
    )
    // TODO If you are using configuration, you need to trouble you to manually modify it
    dbName = "{{$dbName}}"
    user = "{{$user}}"
    password = "{{$password}}"
    host = "{{$host}}"

    dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)

    DB_{{- $dbType }}, err = gorm.Open(mysql.New(mysql.Config{
        DSN:                       dsn,   // DSN data source name
        DefaultStringSize:         256,   // string 类型字段的默认长度
        DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
        DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
        DontSupportRenameColumn:   true,  // 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
        SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
    }), &gorm.Config{
        // TODO Maybe you need to add some configuration
    })
    sqlDB, _ := DB_{{- $dbType }}.DB()
    if err != nil {
        // TODO Or if you don't know how to handle this error
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(10 * time.Minute)
    //DB_{{- $dbType }}.Callback().Create().Before("gorm:create").Register("update_created_at", createdCallBack)
    //DB_{{- $dbType }}.Callback().Update().Before("gorm:update").Register("update_modified_at", updateCallBack)

}

//func createdCallBack(tx *gorm.DB) {
// // TODO Perhaps you want to do something when you create
//}
//
//func updateCallBack(tx *gorm.DB) {
// // TODO Maybe you want to do something when modifying
//}`
	WAREHOUSE_SQLSERVER_TEMPLATE = `{{- $packageName := .Package}}{{ $dbName := .DbName}}{{ $user := .User}}{{ $password := .Password -}}
{{- $host := .Host }}{{ $now := .Now }}{{ $dbType := .DbType -}}
{{"// Package "}}{{$packageName}}{{"This model file is automatically generated in "}}{{$now}}{{" and can be edited at will, and it will not be overwritten when it is generated again"}}
package {{$packageName}}
import (
    "fmt"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
    "time"
)

var DB_{{- $dbType }} *gorm.DB

func Setup() {
    var (
        dbName, user, password, host, dsn string
        err                               error
    )
    // TODO If you are using configuration, you need to trouble you to manually modify it
    dbName = "{{$dbName}}"
    user = "{{$user}}"
    password = "{{$password}}"
    host = "{{$host}}"

    dsn = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",user,password,host,dbName)
    DB_{{- $dbType }}, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    sqlDB, _ := DB_{{- $dbType }}.DB()
    if err != nil {
        // TODO Or if you don't know how to handle this error
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(10 * time.Minute)
    //DB_{{- $dbType }}.Callback().Create().Before("gorm:create").Register("update_created_at", createdCallBack)
    //DB_{{- $dbType }}.Callback().Update().Before("gorm:update").Register("update_modified_at", updateCallBack)

}

//func createdCallBack(tx *gorm.DB) {
// // TODO Perhaps you want to do something when you create
//}
//
//func updateCallBack(tx *gorm.DB) {
// // TODO Maybe you want to do something when modifying
//}`
)
