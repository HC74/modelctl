// Package model_temp This file is generated by Cli, please do not modify
package model_temp


import "time"


type PicEntity struct {
	CreateTime	time.Time	`gorm:"column:CreateTime"`
	CustomerId	int	`gorm:"column:CustomerId"`
	MonthId	int	`gorm:"column:MonthId"`
	FlagDelete	bool	`gorm:"column:FlagDelete"`
	Pic	string	`gorm:"column:Pic"`
	Id	string	`gorm:"column:Id"`
}

// TableName 表名
func (*PicEntity) TableName() string {
    return "PicEntity"
}