package user

import (
	"slingshot/db"
)

// TODO: UserID snowflake
type User struct {
	Id        int64  `json:"id" form:"id" param:"id" xorm:"autoincr index"`
	UserId    int64  `json:"user_id" form:"user_id" param:"user_id" xorm:"index"`
	Username  string `json:"username" form:"username" xorm:"varchar(64) unique notnull default('') index 'username'"`
	Nickname  string `json:"nickname" form:"nickname" xorm:"varchar(64) index default('')"`
	Password  string `json:"password" form:"password" xorm:"varchar(40) default('')"`
	Email     string `json:"email" form:"email" xorm:"varchar(255)"`
	Phone     string `json:"phone" form:"phone" xorm:"varchar(20)"`
	Status    int    `json:"status" form:"status" xorm:"index default(0)"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`

	Roles []Role `json:"roles" xorm:"-"`
}

func (User) TableName() string {
	return "sys_users"
}

// Add user
func (u *User) Add() (int64, error) {
	return db.DB.Insert(u)
}

// Update user
func (u *User) Update() (int64, error) {
	return db.DB.Where("id = ?", u.Id).Update(u)
}

// Delete user
func (u *User) Delete() (int64, error) {
	return db.DB.Where("id = ?", u.Id).Delete(u)
}

// Get user
func (u *User) Get() (bool, error) {
	return db.DB.Where("id = ?", u.Id).Get(u)
}

// Get user's all roles
func (u *User) GetRoles() error {
	return db.DB.Table("sys_roles").Join("INNER", "sys_user_roles", "sys_roles.id = sys_user_roles.role_id").Where("sys_user_roles.user_id = ?", u.Id).Find(&u.Roles)
}

// Get users
func GetUsers(users *[]User) error {
	return db.DB.Find(users)
}

type Role struct {
	Id        int64  `json:"id" form:"id" param:"id" xorm:"autoincr index"`
	Name      string `json:"name" form:"name" query:"name" xorm:"varchar(64) unique notnull default('')"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`

	Users []User `json:"users" xorm:"-"`
}

func (Role) TableName() string {
	return "sys_roles"
}

// Add role
func (r *Role) Add() (int64, error) {
	return db.DB.Insert(r)
}

// Update role
func (r *Role) Update() (int64, error) {
	return db.DB.Where("id = ?", r.Id).Update(r)
}

// Delete role
func (r *Role) Delete() (int64, error) {
	return db.DB.Where("id = ?", r.Id).Delete(r)
}

// Get role
func (r *Role) Get() (bool, error) {
	return db.DB.Where("id = ?", r.Id).Get(r)
}

// Get role's all users
func (r *Role) GetUsers() error {
	return db.DB.Table("sys_users").Join("INNER", "sys_user_roles", "sys_users.id = sys_user_roles.user_id").Where("sys_user_roles.role_id = ?", r.Id).Find(&r.Users)
}

// Get roles
func GetRoles(roles *[]Role) error {
	return db.DB.Find(roles)
}

type UserRole struct {
	ID        int64 `json:"id" form:"id" param:"id" xorm:"autoincr index"`
	UserID    uint  `json:"user_id" form:"user_id" query:"user_id" xorm:"index default 0"`
	RoleID    uint  `json:"role_id" form:"role_id" query:"role_id" xorm:"index default 0"`
	CreatedAt int64 `json:"created_at" xorm:"created"`
	UpdatedAt int64 `json:"updated_at" xorm:"updated"`
	DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
}

func (UserRole) TableName() string {
	return "sys_user_roles"
}

// Add user role
func (ur *UserRole) Add() (int64, error) {
	return db.DB.Insert(ur)
}

// Update user role
func (ur *UserRole) Update() (int64, error) {
	return db.DB.Where("id = ?", ur.ID).Update(ur)
}

// Delete user role
func (ur *UserRole) Delete() (int64, error) {
	return db.DB.Where("id = ?", ur.ID).Delete(ur)
}

// type policy struct
type Policy struct {
	Role   string `json:"role" form:"role" query:"role"`
	Path   string `json:"path" form:"path" query:"path"`
	Method bool   `json:"method" form:"method" query:"method"`
}

// Migrate user
func Migrate() {
	db.DB.Sync(&User{})
	db.DB.Sync(&Role{})
}
