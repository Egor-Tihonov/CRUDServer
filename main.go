package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

type person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Works bool   `json:"works"`
}

var (
	users       = map[int]*person{}
	seq         = 1
	databaseUrl = ""
	conn, err   = pgxpool.Connect(context.Background(), os.Getenv(databaseUrl))
)

func createUser(c echo.Context) error {
	u := &person{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	u := new(person)
	err := c.Bind(u)
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c echo.Context) error {
	var u = new(person)
	row, err2 := conn.Query(context.Background(), "Select * from person")
	if err2 != nil {
		log.Fatal("error")
	}
	for row.Next() {
		u = new(person)
		values, err := row.Values()
		if err != nil {
			log.Fatal("error")
		}
		u.ID = int(values[0].(int32))
		u.Name = values[1].(string)
		u.Works = values[2].(bool)
	}
	return c.JSON(http.StatusOK, u)
}

func main() {

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	e := echo.New()

	// Routes
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	// Start server
	err := e.Start(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
