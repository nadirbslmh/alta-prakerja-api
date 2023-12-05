package main

import (
	"gugcp/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/upload", controllers.UploadFile)

	e.Logger.Fatal(e.Start(":8080"))
}
