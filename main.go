package main

import (
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatal("unable to connect ", err)
	}
	rps := service.NewService(repository.New(pool))
	defer pool.Close()
	h := handlers.NewHandler(rps)
	e := echo.New()

	e.GET("/users", h.GetAllUsers)
	e.POST("/usersCreate", h.CreateUser)
	e.PUT("/usersUpdate/:id", h.UpdateUser)
	e.DELETE("/usersDelete/:id", h.DeleteUser)
	e.GET("/users/:id", h.GetUserById)
	err = e.Start(":8000")
	if err != nil {
		fmt.Println(err)
	}
}
