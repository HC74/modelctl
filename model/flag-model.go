package model

import (
	"github.com/HC74/modelctl/utils"
	"strings"
)

type FlagModel struct {
	// 数据库类型
	DatabaseType string `json:"database_type"`
	// 输出的目录
	OutDir string `json:"out_dir"`
	// 连接字符串
	Url          string   `json:"url"`
	Table        []string `json:"table"`
	TableStr     string   `json:"table_str"`
	PackageName  string   `json:"package_name"`
	DatabaseName string   `json:"database_name"`
}

// InitTables 初始化表
func (f *FlagModel) InitTables() {
	f.Table = strings.Split(f.TableStr, ",")
}

// InitDatabaseName 填充数据库名
func (f *FlagModel) InitDatabaseName() {
	f.DatabaseName, _ = utils.GetDbNameForUrl(f.Url)
}
