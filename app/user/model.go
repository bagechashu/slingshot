package user

import (
	"slingshot/db"
)

type User struct {
	Id        int64  `json:"id" form:"id" param:"id" xorm:"autoincr index"`
	Username  string `json:"username" form:"username" xorm:"varchar(64) unique notnull default('') index 'username'"`
	Nickname  string `json:"nickname" form:"nickname" xorm:"varchar(64) index default('')"`
	Password  string `json:"password" form:"password" xorm:"varchar(40) default('')"`
	Email     string `json:"email" form:"email" xorm:"varchar(255)"`
	Phone     string `json:"phone" form:"phone" xorm:"varchar(20)"`
	Status    int    `json:"status" form:"status" xorm:"index default(0)"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`
}

func (User) TableName() string {
	return "sys_users"
}

type Role struct {
	Name      string `json:"name" form:"name" query:"name" xorm:"varchar(64) unique notnull default('')"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`
}

func (Role) TableName() string {
	return "sys_roles"
}

type UserRole struct {
	UserID    uint  `json:"user_id" form:"user_id" query:"user_id" xorm:"index default 0"`
	RoleID    uint  `json:"role_id" form:"role_id" query:"role_id" xorm:"index default 0"`
	CreatedAt int64 `json:"created_at" xorm:"created"`
	UpdatedAt int64 `json:"updated_at" xorm:"updated"`
	DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
}

func (UserRole) TableName() string {
	return "sys_user_roles"
}

func Migrate() {
	db.DB.Sync(&User{})
	db.DB.Sync(&Role{})
}
