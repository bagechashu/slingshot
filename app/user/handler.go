package user

import (
	"fmt"
	"log"
	"net/http"
	"slingshot/config"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// func login with jwt
func login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := User{Username: req.Username}
	if exist, err := user.GetByUsername(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else if !exist {
		return c.JSON(http.StatusUnauthorized, "username or password error")
	}

	if ok := user.checkUserPassword(req.Password); !ok {
		return c.JSON(http.StatusUnauthorized, "username or password error")
	}

	// return jwt token
	td := fmt.Sprintf("%d%s", config.Cfg.Server.JwtExpiresHour, "h")
	token, err := CreateJwtToken(user.Uid, td)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func register(c echo.Context) error {
	req := RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if exist, err := checkUserEmailExists(req.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else if exist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email already exists",
		})
	}

	user := User{Username: req.Username, Email: req.Email}
	user.setUserPassword(req.Password)

	// Generate short UUID for user id
	id, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to generate user id",
		})
	}
	user.Uid = id.String()[:8]

	if _, err := user.Add(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

func getUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	user.GetByUid()
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
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

func delUser(c echo.Context) error {
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	u, err := user.GetByUid()
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

// func: get role
// param: id
// return: role
func getRole(c echo.Context) error {
	role := Role{}
	if err := c.Bind(&role); err != nil {
		return err
	}
	role.Get()
	return c.JSON(http.StatusOK, role)
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
	// bind request data
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// check if user exist and has no roles
	users := make([]User, 0)
	for _, uid := range requestData.Uids {
		user := User{Uid: uid}
		// check if user exist
		if exist, err := user.GetByUid(); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		} else if !exist {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "User not exist",
			})
		}
		// check if user has roles
		if roles, err := GetRolesOfUser(uid); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		} else if len(roles) > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "User already has roles",
			})
		}
		users = append(users, user)
	}

	// check if role exist
	role := Role{Rid: c.Param("rid")}
	if exist, err := role.Get(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Role not exist",
		})
	}

	// add users to role
	for _, user := range users {
		// update User struct Rid
		user.Rid = role.Rid
		if _, err := user.UpdateRid(); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// casbin add group policy
		if result, err := Rbac.Enforcer.AddGroupingPolicy(user.Uid, role.Rid); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		} else if !result {
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
	policy := Rbac.Enforcer.GetFilteredPolicy(0, role.Name)
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

	if ok, err := policy.IsValidRole(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid role",
		})
	}

	if ok, err := policy.IsValidMethod(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid method",
		})
	}

	Rbac.Enforcer.AddPolicy(policy.Role, policy.Path, policy.Method)

	// update policy
	LoadPolicy()

	return c.JSON(http.StatusOK, policy)
}

// func: Get all policy.
// return: policy
func getPolicys(c echo.Context) error {
	policys := Rbac.Enforcer.GetPolicy()
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

	Rbac.Enforcer.RemovePolicy(policy.Role, policy.Path, policy.Method)

	// update policy
	LoadPolicy()

	return c.JSON(http.StatusOK, policy)
}
