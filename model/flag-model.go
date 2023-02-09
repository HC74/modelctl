package model

import "strings"

type FlagModel struct {
	// 数据库类型
	DatabaseType string `json:"database_type"`
	// 输出的目录
	OutDir string `json:"out_dir"`
	// 连接字符串
	Url         string   `json:"url"`
	Table       []string `json:"table"`
	TableStr    string   `json:"table_str"`
	PackageName string   `json:"package_name"`
}

// InitTables 初始化表
func (f *FlagModel) InitTables() {
	f.Table = strings.Split(f.TableStr, ",")
}
