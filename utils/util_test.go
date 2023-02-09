package utils

import (
	"fmt"
	"modelctl/utils/ioutils"
	"os/user"
	"testing"
)

func TestA(t *testing.T) {
	url := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	url, err := GetDbNameForUrl(url)
	fmt.Println(url)
	if err != nil {

	}
}

func TestPath(t *testing.T) {
	current, _ := user.Current()
	path := fmt.Sprintf("%s/.modelctl", current.HomeDir)
	if ioutils.IsExists(path) {
		fmt.Println("存在")
	} else {
		err := ioutils.Mkdir(path, 0666)
		if err != nil {
			fmt.Println("文件创建失败")
		}
	}
	fmt.Println("end")
}
