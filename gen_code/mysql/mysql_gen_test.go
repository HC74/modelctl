package mysql

import (
	"fmt"
	"modelctl/model"
	"testing"
)

func TestA(t *testing.T) {
	m := &FlagMysqlDbGen{
		flagModel: model.FlagModel{Url: "XX/test"},
		dbGen:     &MysqlDBGen{},
	}
	m.ConnectionDB()
	fmt.Println("链接成功")
	data, err := m.AllTableData()
	if err != nil {

	}
	fmt.Println(data)
}
