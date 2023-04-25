package user

import (
	"log"
	"net/http"
	"slingshot/db"
	mw "slingshot/middleware"

	"github.com/labstack/echo/v4"
)

func getUser(c echo.Context) error {
	id := c.Param("id")
	user := User{}
	db.DB.Where("id = ?", id).Get(&user)
	return c.JSON(http.StatusOK, user)
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
	id := c.Param("id")
	user := User{}
	db.DB.Where("id = ?", id).Find(&user)
	db.DB.Delete(&user)
	return c.JSON(http.StatusOK, user)
}

func addPolicy(c echo.Context) error {
	id := c.Get("id")
	path := c.Get("path")
	method := c.Get("method")
	user := User{}
	db.DB.Where("id = ?", id).Find(&user)

	mw.Rbac.Enforcer.AddPolicy(user.Username, path, method)
	return c.JSON(http.StatusOK, user)
}
