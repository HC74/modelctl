package mysql

import (
	"database/sql"
	"fmt"
	"github.com/HC74/modelctl/model"
	"github.com/HC74/modelctl/utils"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	// 查询所有表
	tableNamesSql   = `select table_name from information_schema.tables where table_schema = ? and table_type = 'BASE TABLE';`
	tableNameColumn = `select column_name,
		is_nullable, if(column_type = 'tinyint(1)', 'boolean', data_type),
		column_type like '%unsigned%'
		from information_schema.columns
		where table_schema = ? and  table_name = ?
		order by ordinal_position;`
)

var flagGen *FlagMysqlDbGen

type MysqlDBGen struct {
	db     *sql.DB
	dbName string
}

type FlagMysqlDbGen struct {
	dbGen     *MysqlDBGen
	flagModel model.FlagModel
}

// NewMysqlDbGen 创建mysql生成器
func NewMysqlDbGen(flag model.FlagModel) *FlagMysqlDbGen {
	return &FlagMysqlDbGen{dbGen: &MysqlDBGen{}, flagModel: flag}
}

// ConnectionDB 测试连接数据库
func (f *FlagMysqlDbGen) ConnectionDB() (err error) {
	flagModel := f.flagModel
	dbName, err := utils.GetDbNameForUrl(flagModel.Url)
	if err != nil {
		return
	}
	f.dbGen.dbName = dbName
	db, err := sql.Open("mysql", flagModel.Url)
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

func (f *FlagMysqlDbGen) TableDataForSelect(flagModel model.FlagModel) (model.TableMataDataList, error) {
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
func (f *FlagMysqlDbGen) AllTableData() (model.TableMataDataList, error) {
	// 获取数据库实例
	db := f.dbGen.db
	// 获取数据库名称
	dbName := f.dbGen.dbName
	rows, err := db.Query(tableNamesSql, dbName)
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

// GetTableColumns 根据表名称查询列元数据集
func (f *FlagMysqlDbGen) GetTableColumns(tableName string) (cols model.ColumnMetaDataList, err error) {
	db := f.dbGen.db
	dbName := f.dbGen.dbName
	rows, err := db.Query(tableNameColumn, dbName, tableName)
	if err != nil {
		return
	}
	cols = model.ColumnMetaDataList{}
	// 一列一列的读取
	for rows.Next() {
		// 名称 是否为空 数据类型
		var name, isNullable, dataType string
		var isUnsigned bool
		if err := rows.Scan(&name, &isNullable, &dataType, &isUnsigned); err != nil {
			return nil, err
		}
		// 是否为空
		isNull := strings.ToLower(isNullable) == "yes"
		cols = append(cols, model.NewColumnMetaData(name, isNull, dataType, isUnsigned, tableName, "", false, "mysql"))
	}
	return cols, rows.Err()
}
