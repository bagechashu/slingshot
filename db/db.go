package db

import (
	"fmt"
	"os"
	"slingshot/config"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
)

var DB *xorm.Engine
var once sync.Once

func InitMysql() {
	// log.Printf("database Dsn: %s", config.Cfg.Database.DSN())
	once.Do(func() {
		var err error
		DB, err = xorm.NewEngine("mysql", config.Cfg.Database.DSN())
		if err != nil {
			fmt.Println("db init failed: ", err.Error())
		}
		DB.SetLogger(xormlog.NewSimpleLogger(os.Stdout))
	})
}
