package user

import (
	"slingshot/db"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" form:"username" query:"username" gorm:"size:64;uniqueIndex;default:'';not null;"`
	Nickname string `json:"nickname" form:"nickname"  gorm:"size:64;index;default:'';"`
	Password string `json:"password" form:"password" gorm:"size:40;default:'';"`
	Email    string `json:"email" form:"email" gorm:"size:255;"`
	Phone    string `json:"phone" form:"phone" gorm:"size:20;"`
	Status   int    `json:"status" form:"status" gorm:"index;default:0;"`
}

func (User) TableName() string {
	return "sys_users"
}

type Role struct {
	gorm.Model
	Name string `json:"name" form:"name" query:"name" gorm:"size:64;uniqueIndex;default:'';not null;"`
}

func (Role) TableName() string {
	return "sys_roles"
}

type UserRole struct {
	gorm.Model
	UserID uint `json:"user_id" form:"user_id" query:"user_id" gorm:"index;default:0;"`
	RoleID uint `json:"role_id" form:"role_id" query:"role_id" gorm:"index;default:0;"`
}

func (UserRole) TableName() string {
	return "sys_user_roles"
}

func Migrate() {
	db.DB.AutoMigrate(&User{})
	db.DB.AutoMigrate(&Role{})
	db.DB.AutoMigrate(&gormadapter.CasbinRule{})
}
