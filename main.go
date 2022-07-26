package main //main

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middleware"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	poolP *pgxpool.Pool
	poolM *mongo.Client
)

func main() {
	dbname := "postgres"
	e := echo.New()
	rdsClient := redisConnection()
	conn := DbConnection(dbname)
	defer func() {
		err := rdsClient.Close()
		if err != nil {
			log.Errorf("error while closing redis connection - %v", err)
		}
		poolP.Close()
		err = poolM.Disconnect(context.Background())
		if err != nil {
			log.Errorf("error close mongo connection - %e", err)
		}
	}()
	c := cache.NewCache(rdsClient)
	rps := service.NewService(conn)
	h := handlers.NewHandler(rps, c)
	e.GET("/users", h.GetAllUsers)
	e.GET("/attachment", h.DownloadFile)
	e.POST("/sign-up", h.Registration)
	e.PUT("/usersUpdate/:id", h.UpdateUser, middleware.IsAuthenticated)
	e.DELETE("/usersDelete/:id", h.DeleteUser, middleware.IsAuthenticated)
	e.POST("/login/:id", h.Authentication)
	e.POST("/logout/:id", h.Logout, middleware.IsAuthenticated)
	e.GET("/users/:id", h.GetUserById, middleware.IsAuthenticated)
	e.GET("/refreshToken", h.RefreshToken, middleware.IsAuthenticated)
	e.POST("/upload", h.Upload)
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
			log.Errorf("bad connection with postgresql: %v", err)
			return nil
		}
		return &repository.PRepository{Pool: poolP}

	case "mongo":
		poolM, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
		if err != nil {
			log.Errorf("bad connection with mongoDb: %v", err)
			return nil
		}
		return &repository.MRepository{Pool: poolM}

	}
	return nil
}
func redisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"status": "connection to redis is failed",
			"err":    err,
		})
		return nil
	}
	log.WithFields(log.Fields{
		"status": "connection with redis was success",
	})
	return rdb
}
