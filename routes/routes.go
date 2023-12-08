package routes

import (
	"gugcp/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())

	apiRoutes := e.Group("/api/v1")

	apiRoutes.POST("/upload", controllers.UploadFile)
}
