package main

import (
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatal("unable to connect ", err)
	}

	rps := repository.New(pool)
	defer pool.Close()
	h := handlers.NewHandler(rps)
	e := echo.New()

	e.GET("/users", h.GetAllUsers)
	e.GET("/usersCreate", h.CreateUser)
	e.GET("/usersUpdate/:id", h.UpdateUser)
	e.GET("/usersDelete/:id", h.DeleteUser)
	e.GET("/users/:id", h.GetUserById)
	err = e.Start(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
