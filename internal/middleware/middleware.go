package middleware

import (
	"awesomeProject/internal/service"
	"github.com/labstack/echo/v4/middleware"
)

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: service.JwtKey,
})
