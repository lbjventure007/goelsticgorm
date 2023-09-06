package model

type Orders struct {
	ID        int64 `gorm:"primarykey"`
	UserID    int64
	ProductID int64
}
