package main //main

import (
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middleware"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	poolP *pgxpool.Pool
	poolM *mongo.Client
)

func main() {
	dbname := "postgres"
	conn := DbConnection(dbname)
	defer func() {
		poolP.Close()
		err := poolM.Disconnect(context.Background())
		if err != nil {
			log.Errorf("error close mongo connection - %e", err)
		}
	}()
	rps := service.NewService(conn)
	h := handlers.NewHandler(rps)
	e := echo.New()
	e.GET("/users", h.GetAllUsers)
	e.GET("/download", h.DownloadImage)
	e.POST("/sign-up", h.Registration)
	e.PUT("/usersUpdate/:id", h.UpdateUser, middleware.IsAuthenticated)
	e.DELETE("/usersDelete/:id", h.DeleteUser, middleware.IsAuthenticated)
	e.POST("/login/:id", h.Authentication)
	e.GET("/users/:id", h.GetUserById, middleware.IsAuthenticated)
	e.GET("/refreshToken", h.RefreshToken, middleware.IsAuthenticated)
	err := e.Start(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

func DbConnection(_dbname string) repository.Repository {
	switch _dbname {
	case "postgres":
		poolP, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
		if err != nil {
			log.Errorf("bad connection with postgresql: ", err)
			return nil
		}
		return &repository.PRepository{Pool: poolP}

	case "mongo":
		poolM, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
		if err != nil {
			log.Errorf("bad connection with mongoDb: ", err)
			return nil
		}
		return &repository.MRepository{Pool: poolM}

	}
	return nil
}
