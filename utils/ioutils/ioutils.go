package ioutils

import (
	"fmt"
	"os"
)

// IsExists 文件/目录 是否存在
func IsExists(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// Mkdir 创建文件夹
func Mkdir(path string, code os.FileMode) error {
	err := os.Mkdir(path, code)
	if err != nil {
		fmt.Println("目录已存在或不能按此目录创建：", path)
		return err
	}
	return nil
}
