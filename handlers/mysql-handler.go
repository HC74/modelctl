package handlers

import (
	"modelctl/model"
)

// MysqlHandle mysql处理器
func MysqlHandle(f *model.FlagModel) {
	// 数据库链接
	_ = f.Url
}
