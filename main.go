package main

import (
	"awesomeProject/handlers"
	"awesomeProject/repository"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	defer repository.Conn.Close()

	e := echo.New()

	e.GET("/users", handlers.GetAllUsers)
	e.GET("/usersCreate", handlers.CreateUser)
	/*e.PUT("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)*/

	err := e.Start(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
