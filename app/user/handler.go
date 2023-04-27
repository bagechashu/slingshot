package user

import (
	"log"
	"net/http"
	mw "slingshot/middleware"

	"github.com/google/uuid"
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
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// Generate short UUID for user id
	id, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to generate user id",
		})
	}
	user.Uid = id.String()[:8]

	// log.Printf("user: %v", user)
	_, err = user.Add()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to add user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func delUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	u, err := user.Get()
	if err != nil {
		return err
	}

	log.Printf("user: %v", u)
	if !u {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User not exist",
		})
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
	// Generate short UUID for user id
	id, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to generate user id",
		})
	}
	role.Rid = id.String()[:4]

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
func getRolesOfUser(c echo.Context) error {
	roles, err := GetRolesOfUser(c.Param("uid"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, roles)
}

// func: Get role's all users
// param: id
// return: users
func getUsersOfRole(c echo.Context) error {
	users, err := GetUsersOfRole(c.Param("rid"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// func: add users for role.
// param: id, users
// return: role
func addUsersForRole(c echo.Context) error {
	requestData := struct {
		Uids []string `json:"uid"`
	}{}
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}
	users := make([]User, 0)
	for _, uid := range requestData.Uids {
		user := User{Uid: uid}
		if exist, err := user.Get(); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		} else if !exist {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "User not exist",
			})
		}
		users = append(users, user)
	}
	// log.Printf("================ users: %v===============\n", users)

	role := Role{Rid: c.Param("rid")}
	if exist, err := role.Get(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Role not exist",
		})
	}
	// log.Printf("================ role: %v===============\n", role)

	for _, user := range users {
		if result, err := mw.Rbac.Enforcer.AddGroupingPolicy(user.Uid, role.Rid); err != nil || !result {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Unable to add users to role",
			})
		}
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
