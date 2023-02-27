package handlers

import (
	"fmt"
	"github.com/HC74/modelctl/model"
	"github.com/HC74/modelctl/templates"
	"github.com/HC74/modelctl/utils"
	"github.com/fatih/color"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

type (
	TmplColStruct struct {
		Column     string // 字段名
		ColumnType string // 字段类型
		Tag        string // 字段的tag
		Remark     string // 备注
	}
	TmplColStructList = []*TmplColStruct
	TmplStruct        struct {
		Package         string            // 包名称
		StructName      string            // 结构体名称
		TableName       string            // 表名称
		HasTime         bool              // 是否包含time类型的字段
		Cols            TmplColStructList // 属性
		HasKey          bool              // 是否包含主键
		KeyColumn       string            // 主键列名
		KeyStructColumn string            // 主键列名对应到结构体的名称
		KeyType         string            // 主键类型
		DBName          string            // DB_数据库类型 组合成仓储通用DB
	}
	TmplWarehouse struct {
		Package  string `json:"package"`  // 包名
		DbName   string `json:"db_name"`  // 数据库名称
		User     string `json:"user"`     // 用户名
		Password string `json:"password"` // 密码
		Host     string `json:"host"`     // 主机
		Now      string `json:"now"`      // 当前时间
		DbType   string `json:"dbType"`   // 数据库类型
	}
	TmplStructList = []*TmplStruct
)

func ParseTemplateHandler(t TmplStructList, flagModel model.FlagModel) int {
	path := fmt.Sprintf("./%v", flagModel.PackageName)
	exists, _ := utils.PathExists(path)
	if !exists {
		err := os.Mkdir(path, 0777)
		if err != nil {
			panic(err)
		}
	}
	return RenderFile(t, flagModel)
}
func RenderFile(t TmplStructList, flagModel model.FlagModel) int {
	var wg sync.WaitGroup
	tpl, _ := template.New("model-tpl").Parse(templates.GetGeneratorStructCode())
	curlTpl, _ := template.New("model-curd-tpl").Parse(templates.GetCURDCode())
	warehouseTpl, _ := template.New("model-warehouse-tpl").Parse(templates.GetWAREHOUSECode(flagModel.DatabaseType))
	wgCount := len(t) * 2
	fmt.Println(wgCount)
	wg.Add(wgCount)
	databaseName, username, password, urlPort := utils.ProcessCdn(flagModel.Url, flagModel.DatabaseType)
	twTpl := &TmplWarehouse{
		Package:  flagModel.PackageName,
		DbName:   databaseName,
		User:     username,
		Password: password,
		Host:     urlPort,
		Now:      time.Now().Format("2006-01-02 15:04:05"),
		DbType:   flagModel.DatabaseType,
	}
	FileHandler(*twTpl, &wg, flagModel.PackageName,
		fmt.Sprintf("%s_warehouse", flagModel.DatabaseName), warehouseTpl, false, false)
	for _, tmplStruct := range t {
		go FileHandler(*tmplStruct, &wg, flagModel.PackageName, tmplStruct.TableName, tpl, true, true)
	}
	for _, tmplStruct := range t {
		go FileHandler(*tmplStruct, &wg, flagModel.PackageName,
			fmt.Sprintf("%s%s", tmplStruct.TableName, "_crud"), curlTpl, false, true)
	}
	wg.Wait()
	return wgCount
}

func FileHandler(t interface{}, wg *sync.WaitGroup, pathName, fileName string, tpl *template.Template, ow, isGo bool) {
	path := fmt.Sprintf("./%v/%v.go", pathName, fileName)
	exists, _ := utils.PathExists(path)
	owe := true
	if exists {
		if !ow {
			color.Red("[SERIOUS_WARN]This %v.go file already exists, skipped", fileName)
			if isGo {
				wg.Done()
			}
			return
		}
		_ = os.Remove(path)
		color.Yellow("[WARN]The %v file already exists and is being overwritten", path)
		owe = false
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	if owe {
		color.Green("[SUC]%v.go Generated successfully", fileName)
	}
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(file, t)
	if isGo {
		wg.Done()
	}
}

// Combination 组合方法
func Combination(tables model.TableMataDataList, flagModel model.FlagModel) TmplStructList {
	var table *model.TableMetaData
	var tmplList TmplStructList
	for i := range tables {
		table = tables[i]
		t := tableHandler(table, flagModel)
		tmplList = append(tmplList, t)
	}
	return tmplList
}

// tableHandler 表处理器
func tableHandler(table *model.TableMetaData, flagModel model.FlagModel) *TmplStruct {
	t := &TmplStruct{
		Package:    flagModel.PackageName,
		StructName: ColumnTableHandler(table.Name),
		TableName:  table.Name,
		HasTime:    getHasTime(table.Columns),
		Cols:       ColHandler(table.Columns),
		DBName:     fmt.Sprintf("DB_%s", flagModel.DatabaseType),
	}
	hasKey, keyName, keyType := getHasKey(table.Columns)
	t.HasKey = hasKey
	t.KeyColumn = keyName
	t.KeyStructColumn = ColumnTableHandler(keyName)
	t.KeyType = keyType
	return t
}

func getHasTime(cols model.ColumnMetaDataList) bool {
	for _, col := range cols {
		if col.GoType == "time.Time" {
			return true
		}
	}
	return false
}

// getHasKey 获取是否有主键并回获取
func getHasKey(cols model.ColumnMetaDataList) (bool, string, string) {
	var column *model.ColumnMetaData
	for i := range cols {
		column = cols[i]
		if column.IsPKey {
			return true, column.Name, column.GoType
		}
	}
	return false, "", ""
}

// ColHandler 列处理器
func ColHandler(cols model.ColumnMetaDataList) TmplColStructList {
	var tcsList TmplColStructList
	//var hasKey = false
	var column *model.ColumnMetaData
	for j := range cols {
		column = cols[j]
		t := &TmplColStruct{
			Column:     ColumnTableHandler(column.Name),
			ColumnType: column.GoType,
		}
		if strings.TrimSpace(column.Remark) != "" {
			t.Remark = fmt.Sprintf(" \t // %v", column.Remark)
		} else {
			t.Remark = ""
		}
		if column.IsPKey {
			t.Tag = fmt.Sprintf("`gorm:\"primaryKey;column:%v\" json:\"%v\"`", column.Name, column.Name)
		} else {
			t.Tag = fmt.Sprintf("`gorm:\"column:%v\" json:\"%v\"`", column.Name, column.Name)
		}
		tcsList = append(tcsList, t)
	}
	return tcsList
}

// ColumnTableHandler 行处理器
func ColumnTableHandler(name string) string {
	if strings.Index(name, "_") == -1 {
		return utils.FirstUpper(name)
	}
	snames := strings.Split(name, "_")
	sb := strings.Builder{}
	for _, sname := range snames {
		sb.WriteString(utils.FirstUpper(sname))
	}
	return sb.String()
}
