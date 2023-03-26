package db

import (
	"fmt"
	"log"
	"slingshot/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() (db *gorm.DB) {
	log.Printf("Connecting to database: %s", config.Cfg.Database.DSN())
	db, err := gorm.Open(mysql.Open(config.Cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection to database failed", err)
		panic(err)
	} else {
		fmt.Println("Connected to database")
	}
	return
}

func DB() *gorm.DB {
	if db == nil {
		db = Init()
	}
	return db
}
