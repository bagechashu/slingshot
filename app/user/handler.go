package user

import (
	"log"
	"net/http"
	"slingshot/db"
	mw "slingshot/middleware"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	user.Get()
	return c.JSON(http.StatusOK, user)
}

func getUsers(c echo.Context) error {
	users := make([]User, 0)
	if err := GetUsers(&users); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func addUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	log.Printf("user: %v", user)
	db.DB.InsertOne(&user)
	return c.JSON(http.StatusOK, user)
}

func delUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	user.Delete()
	return c.JSON(http.StatusOK, user)
}

// func: get all roles
// return: roles
func getRoles(c echo.Context) error {
	roles := make([]Role, 0)
	GetRoles(&roles)
	return c.JSON(http.StatusOK, roles)
}

// func: add roles
// param: role
// return: role
func addRole(c echo.Context) error {
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	role.Add()
	return c.JSON(http.StatusOK, role)
}

// func: delete role
// param: id
// return: role
func delRole(c echo.Context) error {
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	role.Delete()
	return c.JSON(http.StatusOK, role)
}

// func: Get user's all roles
// param: id
// return: roles
func getUserRoles(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	user.GetRoles()
	return c.JSON(http.StatusOK, user.Roles)
}

// func: Get role's all users
// param: id
// return: users
func getRoleUsers(c echo.Context) error {
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	role.GetUsers()
	return c.JSON(http.StatusOK, role.Users)
}

// func: add users for role.
// param: id, users
// return: role
func addUsersForRole(c echo.Context) error {
	users := make([]User, 0)
	if err := c.Bind(&users); err != nil {
		return err
	}
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	for _, user := range users {
		mw.Rbac.Enforcer.AddGroupingPolicy(user, role.Name)
	}
	return c.JSON(http.StatusOK, role)
}

// func: Get role's all policy.
// param: id
// return: policy
func getRolePolicy(c echo.Context) error {
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	role.Get()
	policy := mw.Rbac.Enforcer.GetFilteredPolicy(0, role.Name)
	return c.JSON(http.StatusOK, policy)
}

// func: Add Policy for role.
// param: role, path, method
// return: user
func addPolicyForRole(c echo.Context) error {
	policy := Policy{}
	if err := c.Bind(&policy); err != nil {
		return err
	}

	mw.Rbac.Enforcer.AddPolicy(policy.Role, policy.Path, policy.Method)
	return c.JSON(http.StatusOK, policy)
}

// func: Get all policy.
// return: policy
func getPolicys(c echo.Context) error {
	policys := mw.Rbac.Enforcer.GetPolicy()
	return c.JSON(http.StatusOK, policys)
}

// func: Delete Policy for role.
// param: role, path, method
// return: user
func delPolicyForRole(c echo.Context) error {
	policy := Policy{}
	if err := c.Bind(&policy); err != nil {
		return err
	}

	mw.Rbac.Enforcer.RemovePolicy(policy.Role, policy.Path, policy.Method)
	return c.JSON(http.StatusOK, policy)
}
