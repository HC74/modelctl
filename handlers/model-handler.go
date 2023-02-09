package handlers

import (
	"fmt"
	"modelctl/model"
	"modelctl/templates"
	"modelctl/utils"
	"os"
	"strings"
	"sync"
	"text/template"
)

type (
	TmplColStruct struct {
		Column     string // 字段名
		ColumnType string // 字段类型
		Tag        string // 字段的tag
	}
	TmplColStructList = []*TmplColStruct
	TmplStruct        struct {
		Package    string            // 包名称
		StructName string            // 结构体名称
		TableName  string            // 表名称
		HasTime    bool              // 是否包含time类型的字段
		Cols       TmplColStructList // 属性
		Remark     string            // 备注
	}
	TmplStructList = []*TmplStruct
)

var wg sync.WaitGroup

func ParseTemplateHandler(t TmplStructList, flagModel model.FlagModel) int {
	tpl, err := template.New("model-tpl").Parse(templates.GetGeneratorStructCode())
	path := fmt.Sprintf("./%v", flagModel.PackageName)
	exists, _ := utils.PathExists(path)
	if !exists {
		err = os.Mkdir(path, 0777)
		if err != nil {
			panic(err)
		}
	}
	wg.Add(len(t))
	for _, tmplStruct := range t {
		go FileHandler(*tmplStruct, flagModel.PackageName, tpl)
	}
	wg.Wait()
	return len(t)
}

func FileHandler(t TmplStruct, packageName string, tpl *template.Template) {
	path := fmt.Sprintf("./%v/%v.go", packageName, t.TableName)
	exists, _ := utils.PathExists(path)
	if exists {
		_ = os.Remove(path)
		fmt.Println(fmt.Sprintf("The %v file already exists and is being overwritten", path))
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	fmt.Println(fmt.Sprintf("%v.go Generated successfully", t.TableName))
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(file, t)
	wg.Done()
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
	}
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

// ColHandler 列处理器
func ColHandler(cols model.ColumnMetaDataList) TmplColStructList {
	var tcsList TmplColStructList
	var column *model.ColumnMetaData
	for j := range cols {
		column = cols[j]
		t := &TmplColStruct{
			Column:     ColumnTableHandler(column.Name),
			ColumnType: column.GoType,
			Tag:        fmt.Sprintf("`gorm:\"column:%v\"`", column.Name),
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
