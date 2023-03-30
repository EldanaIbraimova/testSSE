package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sse-test/model"
	"time"
)

var dbase *gorm.DB

func Init() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=sse-test port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Message{})
	return db
}

func GetDB() *gorm.DB {
	if dbase == nil {
		dbase = Init()
		var sleep = time.Duration(1)
		for dbase == nil {
			sleep = sleep * 2
			fmt.Printf("Database is unavailable. Wait for %d sec. \n", sleep)
			time.Sleep(sleep * time.Second)
			dbase = Init()
		}
	}
	return dbase
}
