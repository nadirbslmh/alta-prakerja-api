package controllers

import (
	"gugcp/models"
	"gugcp/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GenerateURL(c echo.Context) error {
	var input models.GenerateInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "invalid request",
		})
	}

	services.GenerateURL(input)

	return c.JSON(http.StatusOK, models.Response[any]{
		Status:  true,
		Message: "url generated!",
	})
}
