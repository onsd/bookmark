package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/onsd/bookmark/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// initialize echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Authorized route
	e.POST("/auth", controller.Auth)
	//Login route
	e.POST("/login", controller.Login)

	//Unauthenticated route
	e.GET("/", controller.Accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte(os.Getenv("SECRETKEY"))))
	r.GET("", controller.Restricted)

	e.Start(":" + os.Getenv("PORT"))
}
