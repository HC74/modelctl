package main

import (
	"flag"
	"fmt"
	"modelctl/gen_code"
	"modelctl/handlers"
	"modelctl/model"
	"modelctl/utils"
	"strings"
)

var (
	// 接受命令行输入的参数
	flagModel = model.FlagModel{}
)

func init() {
	flag.StringVar(&flagModel.DatabaseType, "t", "mysql", "数据库类型: 默认为mysql sqlserver:mssql")
	flag.StringVar(&flagModel.Url, "url", "NULL", "数据库链接 默认为空 例如 mysql: 账号:密码@tcp(IP:端口)/库 / sqlserver: server=IP:端口;database=库;userId=账号;password=密码")
	flag.StringVar(&flagModel.TableStr, "tables", "*", "表: 默认为全部,如果要选表生成 则需要 表1,表2,表3....")
	flag.StringVar(&flagModel.OutDir, "dir", "/", "生成目录: 默认为当前目录 暂不支持自定义")
	flag.StringVar(&flagModel.PackageName, "f", "model", "包名 默认为model")
}

func main() {
	// 数据库类型
	flag.Parse()
	// 给Tables赋值
	flagModel.InitTables()
	if flagModel.Url == "NULL" {
		panic("URL不能为空")
	}
	if flagModel.OutDir != "/" {
		panic("暂不支持自定义输出目录")
	}
	// 创建对应的代码生成器
	dbMataData, err := gen_code.NewDbCodeGen(flagModel)
	if utils.IsEmpty(flagModel.PackageName) {
		flagModel.PackageName = "model"
	}
	if err != nil {
		fmt.Println("类型异常，请检查类型是否存在 -type=?值有误")
	}
	if err := dbMataData.ConnectionDB(); err != nil {
		fmt.Println(fmt.Sprintf("ERROR URL : %v", flagModel.Url))
		fmt.Printf("数据库连接时触发恐慌,%s \n", err.Error())
	}
	tables := model.TableMataDataList{}
	// 如果表名未输入或者输入的为* 则默认为所有的表
	if len(flagModel.Table) == 0 || strings.EqualFold(flagModel.TableStr, "*") {
		tables, err = dbMataData.AllTableData()
	} else {
		tables, err = dbMataData.TableDataForSelect(flagModel)
	}
	combination := handlers.Combination(tables, flagModel)
	fileNum := handlers.ParseTemplateHandler(combination, flagModel)
	fmt.Println(fmt.Sprintf("Successfully generated, a total of %v files", fileNum))
	if err != nil {
		fmt.Printf("get meta tables fail %s /n", err.Error())
		return
	}
}
