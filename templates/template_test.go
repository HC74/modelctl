package templates

import (
	"fmt"
	"github.com/HC74/modelctl/utils/ioutils"
	"html/template"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	//package_name := "template"
	m := map[string]string{
		"package_name": "template",
		"dns":          "XX?charset=utf8mb4&parseTime=True&loc=Local",
	}
	tepl, err := template.New("tableTemplate").Funcs(template.FuncMap{
		"now": func() string {
			return time.Now().Format(time.RFC3339)
		},
	}).Parse(GetWarehouseCode())
	if err != nil {
		fmt.Printf("发生了异常 %s \n", err.Error())
	}
	filename := "model.go"
	//caller := getCurrentAbPathByCaller()
	path := "./model"
	var file *os.File
	filePath := fmt.Sprintf("%s/%s", path, filename)
	if !ioutils.IsExists(path) {
		err := ioutils.Mkdir(path, 0666)
		if err != nil {
			fmt.Println("创建文件失败")
		}
	}
	fmt.Println(filePath)
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("创建失败" + err.Error())
	}
	err = os.Chmod(filePath, 0777)
	if err != nil {
		fmt.Println("添加权限失败")
	}
	err = tepl.Execute(file, m)
	if err != nil {
		fmt.Printf("在转换中发生了异常 %s \n", err.Error())
	}
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
