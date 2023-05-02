package config

import (
	"fmt"
)

type Database struct {
	User     string `validate:"required"`
	Password string `validate:"required"`
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Name     string `validate:"required"`
	Param    string `validate:"required"`
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
