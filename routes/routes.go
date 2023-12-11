package routes

import (
	"gugcp/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}" + "\n",
	}))

	apiRoutes := e.Group("/api/v1")

	// generate
	apiRoutes.POST("/url/generate", controllers.GenerateURL)

	// redeem
	apiRoutes.POST("/redeem/save", controllers.SaveRedeemCode)
	apiRoutes.POST("/redeem/check", controllers.CheckStatus)

	// upload task
	apiRoutes.POST("/task/upload", controllers.UploadFile)
}
