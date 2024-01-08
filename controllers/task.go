package controllers

import (
	"gugcp/models"
	"gugcp/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SendFeedback(c echo.Context) error {
	taskID := c.Param("taskID")

	tID, err := strconv.Atoi(taskID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "task ID is invalid",
		})
	}

	var input models.FeedbackInput

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

	result, err := services.SendFeedback(c.Request().Context(), input, tID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[models.FeedbackResponse]{
		Status:  true,
		Message: "feedback sent successfully",
		Data:    result,
	})
}
