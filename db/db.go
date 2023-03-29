package db

import (
	"fmt"
	"slingshot/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() {
	// log.Printf("database Dsn: %s", config.Cfg.Database.DSN())
	var err error
	DB, err = gorm.Open(mysql.Open(config.Cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		fmt.Println("db init failed: ", err.Error())
	}
}
