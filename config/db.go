package config

import (
	"fmt"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	Param    string
}

func (db Database) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Param)
}

func (db Database) DSNWithoutAccount() string {
	return fmt.Sprintf("tcp(%s:%d)/%s?%s",
		db.Host,
		db.Port,
		db.Name,
		db.Param)
}
