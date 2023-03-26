package user

import (
	"log"
	"net/http"
	"slingshot/db"

	"github.com/labstack/echo/v4"
)

func getUsers(c echo.Context) error {
	db := db.DB()
	users := []User{}
	db.Find(&users)
	// spew.Dump(json.Marshal(users))
	// return c.JSON(http.StatusOK, users)
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	db := db.DB()
	id := c.Param("id")
	user := User{}
	db.First(&user, id)
	return c.JSON(http.StatusOK, user)
}

func addUser(c echo.Context) error {
	db := db.DB()
	user := User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	log.Printf("user: %v", user)
	db.Create(&user)
	return c.JSON(http.StatusOK, user)
}

func delUser(c echo.Context) error {
	db := db.DB()
	id := c.Param("id")
	user := User{}
	db.First(&user, id)
	db.Delete(&user)
	return c.JSON(http.StatusOK, user)
}
