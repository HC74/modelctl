package utils

import (
	"errors"
	"fmt"
	"github.com/HC74/modelctl/utils/ioutils"
	"html/template"
	"net/url"
	"os"
	"strings"
)

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

//PathExists 判断一个文件或文件夹是否存在
//输入文件路径，根据返回的bool值来判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func IndexOf(data []string, v string) int {
	for i, item := range data {
		if item == v {
			return i
		}
	}
	return -1
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func GetDbNameForUrl(dns string) (dbName string, err error) {
	index := strings.LastIndex(dns, "/")
	if index <= 0 {
		return "", errors.New("url 输入有误")
	}
	dbNameUrl := dns[index:]
	dbNameUrlObj, err := url.Parse(dbNameUrl)
	if err != nil {
		return
	}
	return strings.Trim(dbNameUrlObj.Path, "/"), nil
}

// ProcessCdn 处理cdn
func ProcessCdn(cdn, t string) (string, string, string, string) {
	if t == "mssql" {
		return ProcessMssqlDBCdn(cdn)
	}
	return ProcessMysqlDBCdn(cdn)
}

func ProcessMssqlDBCdn(dns string) (string, string, string, string) {
	var databaseName, username, password string
	var m = make(map[string]string)
	for _, item := range strings.Split(dns, ";") {
		kv := strings.Split(item, "=")
		key := kv[0]
		val := kv[1]
		m[key] = val
	}
	databaseName = m["database"]
	username = m["userId"]
	password = m["password"]
	return databaseName, username, password, m["server"]
}

func ProcessMysqlDBCdn(cdn string) (string, string, string, string) {
	var databaseName, username, password string
	// 分割掉? 后面的部分
	if strings.Index(cdn, "?") != -1 {
		cdn = SplitGetFirst(cdn, "?")
	}
	// 处理库名
	if strings.Index(cdn, "/") == -1 {
		panic("未找到库名")
	}
	databaseName = SplitGetLast(cdn, "/")
	cdn = SplitGetFirst(cdn, "/")
	if strings.Index(cdn, "@tcp") == -1 {
		panic("未找到url")
	}
	userpwd := SplitGetFirst(cdn, "@tcp")
	urlPort := SplitGetLast(cdn, "@tcp")
	username = SplitGetFirst(userpwd, ":")
	password = SplitGetLast(userpwd, ":")
	urlPort = urlPort[1:]
	urlPort = urlPort[:len(urlPort)-1]
	return databaseName, username, password, urlPort
}

// SplitGetFirst 分割获取第一个元素
func SplitGetFirst(s, sub string) string {
	if strings.Index(s, sub) == -1 {
		return ""
	}
	ss := strings.Split(s, sub)
	return ss[0]
}

// SplitGetLast 分割获取最后一个元素
func SplitGetLast(s, sub string) string {
	if strings.Index(s, sub) == -1 {
		return ""
	}
	ss := strings.Split(s, sub)
	return ss[len(ss)-1]
}

// Template

// NewTemplate 创建模板
// param name 模板名称
// param content 模板内容
// param savePath 保存的路径
// param filename 文件名称
// param map placeholder 占位符 { "key":"value" } => {{ .key }}
func NewTemplate(name, content, savePath, filename string, placeholder map[interface{}]interface{}) {
	// 创建模板
	tepl, err := template.New(name).Parse(content)
	if err != nil {
		fmt.Printf("发生了异常 %s \n", err.Error())
	}
	var file *os.File
	// 拼接为最终生成文件的目录
	filePath := fmt.Sprintf("%s/%s", savePath, fmt.Sprintf("%s.go", filename))
	// 判断目录是否存在 如果不存在则创建
	if !ioutils.IsExists(savePath) {
		err := ioutils.Mkdir(savePath, 0666)
		if err != nil {
			fmt.Println("创建文件失败")
		}
	}
	// 直接覆盖文件
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("创建失败" + err.Error())
	}
	// 渲染模板 写入文件
	err = tepl.Execute(file, placeholder)
	if err != nil {
		fmt.Printf("在转换中发生了异常 %s \n", err.Error())
	}
}
