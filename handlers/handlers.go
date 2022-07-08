package handlers

import (
	"awesomeProject/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUser(c echo.Context) error {
	err := repository.Create()
	return c.JSON(http.StatusOK, err)
}

/*
func UpdateUser(c echo.Context) error {
	u := new(repository.Person)
	err := c.Bind(u)
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}*/

func GetAllUsers(c echo.Context) error {
	repository.SelectAll()
	return c.JSON(http.StatusOK, repository.People)
}
