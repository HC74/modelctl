package gen_code

import (
	"errors"
	"fmt"
	"github.com/HC74/modelctl/gen_code/mssql"
	"github.com/HC74/modelctl/gen_code/mysql"
	"github.com/HC74/modelctl/model"
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
