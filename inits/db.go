package inits

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var once2 sync.Once
var Db *gorm.DB

func InitDb() {
	once2.Do(func() {
		dsn := "root:1234qwer@tcp(127.0.0.1:3307)/sharding_db?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		Db = db
	})

}


