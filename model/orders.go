package model

import "time"

type Orders struct {
	ID      int64 `gorm:"primarykey"`
	UserID  int64
	OrderId int64
	//	ProductID  int64
	Name       string
	CreateTime time.Time
}

func (Orders) TableName() string {
	return "t_order"
}
