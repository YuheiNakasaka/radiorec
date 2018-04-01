package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/YuheiNakasaka/radiorec/internal/handler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/programs", handler.Programs)

	e.Logger.Fatal(e.Start(":1323"))
}
