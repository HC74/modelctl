package templates

import (
	"modelctl/model"
	"modelctl/utils"
)

type ITemplate interface {
	// SetPath 设置模板路径
	SetPath(path string) error
	// HandlerTemplate 处理模板
	HandlerTemplate() error
}

type Template struct {
	template ITemplate
}

// NewTemplate 创建模板
func NewTemplate(f model.FlagModel) *Template {
	return &Template{template: &defaultTemplate{flag: f}}
}

// NewTemplateForPath 待开发
func NewTemplateForPath(string) {

}

type defaultTemplate struct {
	flag model.FlagModel
}

func (d defaultTemplate) SetPath(path string) error {
	return nil
}
func (d defaultTemplate) HandlerTemplate() error {
	flagModel := d.flag
	package_name := flagModel.PackageName
	out_dir := flagModel.OutDir
	// 处理仓储
	utils.NewTemplate("defaultRepostry", GetWarehouseCode(), out_dir, "model", map[interface{}]interface{}{
		"package_name": package_name,
	})
	// 处理具体的实体

	return nil
}
