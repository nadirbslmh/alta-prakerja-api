package routes

import (
	"gugcp/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 		AllowOrigins: []string{"https://one.alterra.academy"},
		AllowOrigins: []string{"*"}, //TODO: change to the origin!
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodPost},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}" + "\n",
	}))

	apiRoutes := e.Group("/api/v1")

	// generate
	apiRoutes.POST("/url/generate", controllers.GenerateURL)

	// redeem
	apiRoutes.POST("/redeem/save", controllers.SaveRedeemCode)
	apiRoutes.POST("/redeem/get/:state", controllers.GetRedeem)
	apiRoutes.POST("/redeem/user/:userID", controllers.GetRedeemByUserID)
	apiRoutes.POST("/redeem/check", controllers.CheckStatus)

	// upload task
	apiRoutes.POST("/task/upload", controllers.UploadFile)

	// task management
	apiRoutes.POST("/task/feedback/:taskID", controllers.SendFeedback)
	apiRoutes.POST("/tasks", controllers.GetAllTasks)
}
