package user

import (
	"errors"
	"regexp"
	"slingshot/db"

	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64  `json:"id" form:"id" param:"id"` // can't add xorm tag, because it's primary key
	Uid       string `json:"uid" form:"uid" param:"uid" xorm:"varchar(8) index"`
	Rid       string `json:"rid" form:"rid" xorm:"varchar(6) index default('')"`
	Username  string `json:"username" form:"username" xorm:"varchar(20) default('') index 'username'"`
	Nickname  string `json:"nickname" form:"nickname" xorm:"varchar(20) index default('')"`
	Password  string `json:"password" form:"password" xorm:"varchar(72) default('')"`
	Email     string `json:"email" form:"email" xorm:"varchar(100) default('') unique notnull"`
	Phone     string `json:"phone" form:"phone" xorm:"varchar(20)"`
	Status    int    `json:"status" form:"status" xorm:"index default(0)"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required"`
	Password string `json:"password" form:"password" query:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"required"`
	Email    string `json:"email" form:"email" query:"email" validate:"required,email"`
	Password string `json:"password" form:"password" query:"password" validate:"required"`
}

func (User) TableName() string {
	return "sys_users"
}

// New user
func NewUser() *User {
	return new(User)
}

// check user email exists
func checkUserEmailExists(email string) (bool, error) {
	if email == "" {
		return false, errors.New("email is empty")
	}

	user := User{Email: email}
	return db.DB.Get(&user)
}

func (u *User) setUserPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if len(hashedPassword) == 0 {
		return errors.New("password hash is empty")
	}

	u.Password = cast.ToString(hashedPassword)
	return nil
}

func (u *User) checkUserPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Add user
func (u *User) Add() (int64, error) {
	return db.DB.Insert(u)
}

// Update user
func (u *User) Update() (int64, error) {
	return db.DB.Where("uid = ?", u.Uid).Update(u)
}

// Update user Rid
func (u *User) UpdateRid() (int64, error) {
	return db.DB.Where("uid = ?", u.Uid).Cols("rid").Update(u)
}

// Delete user
func (u *User) Delete() (int64, error) {
	return db.DB.Where("uid = ?", u.Uid).Delete(u)
}

func (u *User) GetByUsername() (bool, error) {
	return db.DB.Where("username = ?", u.Username).Get(u)
}

func (u *User) GetByUid() (bool, error) {
	return db.DB.Where("uid = ?", u.Uid).Get(u)
}

// Get role of user
func (u *User) GetUidByRid() (bool, error) {
	return db.DB.Where("uid = ?", u.Uid).Cols("rid").Get(u)
}

// Get users
func GetUsers(users *[]User) error {
	return db.DB.Find(users)
}

type Role struct {
	Id        int64  `json:"id" form:"id" param:"id"`
	Rid       string `json:"rid" form:"rid" param:"rid" xorm:"varchar(6) index"`
	Name      string `json:"name" form:"name" query:"name" xorm:"varchar(20) unique notnull default('')"`
	CreatedAt int64  `json:"created_at" xorm:"created"`
	UpdatedAt int64  `json:"updated_at" xorm:"updated"`
	DeletedAt int64  `json:"deleted_at" xorm:"deleted"`
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
	return db.DB.Where("rid = ?", r.Rid).Update(r)
}

// Delete role
func (r *Role) Delete() (int64, error) {
	return db.DB.Where("rid = ?", r.Rid).Delete(r)
}

// Get role
func (r *Role) Get() (bool, error) {
	return db.DB.Where("rid = ?", r.Rid).Get(r)
}

// Get roles
func GetRoles(roles *[]Role) error {
	return db.DB.Find(roles)
}

// TODO: use User and Role model to replace this
// Get roles of user
func GetRolesOfUser(uid string) (roles []Role, err error) {
	rid := make([]string, 0)
	err = db.DB.Table("sys_casbin_rule").Where("p_type = ?", "g").Where("v0 = ?", uid).Cols("v1").Find(&rid)
	if err != nil {
		return nil, err
	}

	for _, v := range rid {
		role := Role{Rid: v}
		_, err := role.Get()
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// TODO: use User and Role model to replace this
// Get users of role
func GetUsersOfRole(rid string) (users []User, err error) {
	uid := make([]string, 0)
	err = db.DB.Table("sys_casbin_rule").Where("p_type = ?", "g").Where("v1 = ?", rid).Cols("v0").Find(&uid)
	if err != nil {
		return nil, err
	}

	for _, v := range uid {
		user := User{Uid: v}
		_, err := user.GetByUid()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// type policy struct
type Policy struct {
	Role   string `json:"role" form:"role" query:"role"`
	Path   string `json:"path" form:"path" query:"path"`
	Method string `json:"method" form:"method" query:"method"`
}

// is valid role
func (p *Policy) IsValidRole() (bool, error) {
	role := Role{Rid: p.Role}
	return role.Get()
}

// is valid method
func (p *Policy) IsValidMethod() (bool, error) {
	pattern := regexp.MustCompile("^GET|POST|PUT|PATCH|DELETE$")
	if !pattern.MatchString(p.Method) {
		return false, errors.New("invalid Request Method")
	}
	return true, nil
}

// Migrate user
func Migrate() {
	db.DB.Sync2(&User{})
	db.DB.Sync2(&Role{})
}
