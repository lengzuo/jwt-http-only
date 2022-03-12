package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"jwt-http-only/backend/vendor/github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/login", Login)
	e.GET("/user", UserAPI)
	e.POST("/refresh", RefreshToken)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
