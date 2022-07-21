package middleware

import (
	"awesomeProject/internal/service"
	"github.com/labstack/echo/v4/middleware"
)

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{ //check user is he Authenticated
	SigningKey: service.JwtKey,
})
