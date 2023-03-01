package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"slingshot/config"
)

var (
	db *gorm.DB
)

func generateDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.Name)
}

func DBInit() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open(generateDSN()), &gorm.Config{})
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
		db = DBInit()
	}
	return db
}
