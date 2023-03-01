package user

import (
	"slingshot/db"

	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	gorm.Model
}

func Migrate() {
	db := db.DB()
	db.AutoMigrate(&User{})
}
