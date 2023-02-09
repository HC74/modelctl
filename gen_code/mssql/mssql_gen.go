package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/HC74/modelctl/model"
	"github.com/HC74/modelctl/utils"
	"strings"
)
import _ "github.com/denisenkom/go-mssqldb"

const (
	// 查询所有表
	tableNamesSql   = `select name from sysobjects where xtype='u'`
	tableNameColumn = `select 
	 b.name column_name,
	 b.IsNullable is_nullable,
	 c.name  data_type
	from sysobjects a,syscolumns b,systypes c where a.id=b.id
	and a.name= ? and a.xtype='U'
	and b.xtype=c.xtype`
	tableGetObjectsId       = `select id from sysobjects a where a.xtype='U' and a.name = ?`
	tableGetColumnMap       = `select name,colid from syscolumns where id = ?`
	tableGetColumnRemarkMap = `select minor_id,[value] from sys.extended_properties where major_id = ?`
	tableGetKey             = `select column_name from information_schema.key_column_usage where table_name = ?`
)

var flagGen *FlagMsSqlDbGen

type MsSqlDBGen struct {
	db     *sql.DB
	dbName string
}

type FlagMsSqlDbGen struct {
	dbGen     *MsSqlDBGen
	flagModel model.FlagModel
}

// NewMsSqlDbGen 创建sqlserver生成器
func NewMsSqlDbGen(flag model.FlagModel) *FlagMsSqlDbGen {
	flag.Url = strings.ReplaceAll(flag.Url, "userId", "user id")
	if strings.Index(flag.Url, ";encrypt=disable") == -1 {
		if !strings.HasSuffix(flag.Url, ";") {
			flag.Url += ";"
		}
		flag.Url += "encrypt=disable"
	}
	return &FlagMsSqlDbGen{dbGen: &MsSqlDBGen{}, flagModel: flag}
}

// ConnectionDB 测试连接数据库
func (f *FlagMsSqlDbGen) ConnectionDB() (err error) {
	flagModel := f.flagModel
	urls := strings.Split(flagModel.Url, ";")
	for _, url := range urls {
		if strings.Index(url, "=") == -1 {
			return errors.New("URL异常")
		}
		if strings.Index(url, "database") == -1 {
			continue
		}
		dbName := strings.Split(url, "=")[1]
		if strings.TrimSpace(dbName) == "" {
			return errors.New("数据库名称未填写")
		}
		f.dbGen.dbName = dbName
	}
	db, err := sql.Open("mssql", flagModel.Url)
	//defer db.Close()
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	f.dbGen.db = db
	fmt.Printf("mysql connection %s", flagModel.Url)
	return
}

func (f *FlagMsSqlDbGen) TableDataForSelect(flagModel model.FlagModel) (model.TableMataDataList, error) {
	// 获取数据库实例
	db := f.dbGen.db
	rows, err := db.Query(tableNamesSql)
	if err != nil {
		return nil, err
	}
	var tables = flagModel.Table
	defer rows.Close()
	var tableMataDataList model.TableMataDataList
	for rows.Next() {
		tableMeta := &model.TableMetaData{}
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		if utils.IndexOf(tables, tableName) == -1 {
			continue
		}
		tableMeta.Name = tableName
		cols, err := f.GetTableColumns(tableName)
		if err != nil {
			return nil, err
		}
		tableMeta.Columns = cols
		tableMataDataList = append(tableMataDataList, tableMeta)
	}
	return tableMataDataList, err
}

// AllTableData 所有表的数据
func (f *FlagMsSqlDbGen) AllTableData() (model.TableMataDataList, error) {
	// 获取数据库实例
	db := f.dbGen.db
	rows, err := db.Query(tableNamesSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tableMataDataList model.TableMataDataList
	for rows.Next() {
		tableMeta := &model.TableMetaData{}
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableMeta.Name = tableName
		cols, err := f.GetTableColumns(tableName)
		if err != nil {
			return nil, err
		}
		tableMeta.Columns = cols
		tableMataDataList = append(tableMataDataList, tableMeta)
	}
	return tableMataDataList, err
}

// SearchKey 查找主键
func SearchKey(db *sql.DB, tableName string) map[string]bool {
	rows, _ := db.Query(tableGetKey, tableName)
	m := make(map[string]bool)
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		m[name] = true
	}
	return m
}

// SearchRemark 查找备注
func SearchRemark(db *sql.DB, tableName string) map[string]string {
	// 获取ID
	id := GetObjectID(db, tableName)
	// 获取 Name 和 colid
	m := GetColumnNameAndColid(db, id)
	// 获取 colid 和 备注
	remarkMap := GetColumnColidAndRemark(db, id, m)
	return remarkMap
}

// GetObjectID 获取系统 objects 表中的id
func GetObjectID(db *sql.DB, tableName string) string {
	rows, _ := db.Query(tableGetObjectsId, tableName)
	var id string
	for rows.Next() {
		_ = rows.Scan(&id)
	}
	return id
}

// GetColumnColidAndRemark 获取colid ID + 备注
func GetColumnColidAndRemark(db *sql.DB, id string, m map[string]string) map[string]string {
	rows, _ := db.Query(tableGetColumnRemarkMap, id)
	rm := make(map[string]string)
	for rows.Next() {
		var colidId, remark string
		_ = rows.Scan(&colidId, &remark)
		if v := m[colidId]; v != "" {
			rm[v] = remark
		}
	}
	return rm
}

// GetColumnNameAndColid 获取列的名称+colid ID
func GetColumnNameAndColid(db *sql.DB, id string) map[string]string {
	rows, _ := db.Query(tableGetColumnMap, id)
	m := make(map[string]string)
	for rows.Next() {
		var name, colid string
		_ = rows.Scan(&name, &colid)
		m[colid] = name
	}
	return m
}

// GetTableColumns 根据表名称查询列元数据集
func (f *FlagMsSqlDbGen) GetTableColumns(tableName string) (cols model.ColumnMetaDataList, err error) {
	db := f.dbGen.db
	// 查询主键和备注
	remarkMap := SearchRemark(db, tableName)
	km := SearchKey(db, tableName)
	rows, err := db.Query(tableNameColumn, tableName)
	if err != nil {
		return
	}
	cols = model.ColumnMetaDataList{}
	// 一列一列的读取
	for rows.Next() {
		// 名称 是否为空 数据类型
		var name, isNullable, dataType string
		var isUnsigned bool
		if err := rows.Scan(&name, &isNullable, &dataType); err != nil {
			return nil, err
		}
		// 是否为空
		isNull := isNullable == "1"
		cols = append(cols, model.NewColumnMetaData(name, isNull, dataType, isUnsigned, tableName, remarkMap[name], km[name], "mssql"))
	}
	return cols, rows.Err()
}
