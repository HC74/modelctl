package gen_code

import (
	"errors"
	"fmt"
	"modelctl/gen_code/mssql"
	"modelctl/gen_code/mysql"
	"modelctl/model"
	"strings"
)

// NewDbCodeGen 新建数据库生成器
func NewDbCodeGen(f model.FlagModel) (IDBMetaData, error) {
	switch strings.ToLower(f.DatabaseType) {
	case "mysql":
		return mysql.NewMysqlDbGen(f), nil
	case "mssql":
		return mssql.NewMsSqlDbGen(f), nil
	}
	return nil, errors.New(fmt.Sprintf("未知的类型:%s", f.DatabaseType))
}
