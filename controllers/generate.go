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

	if err := input.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response[[]*models.ValidationErrorResponse]{
			Status:  false,
			Message: "request validation failed",
			Data:    err,
		})
	}

	response, err := services.GenerateURL(input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[models.Data]{
		Status:  true,
		Message: "url generated successfully",
		Data:    response.Data,
	})
}
