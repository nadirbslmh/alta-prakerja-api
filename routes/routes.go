package routes

import (
	"gugcp/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	e.POST("/upload", controllers.UploadFile)
}
