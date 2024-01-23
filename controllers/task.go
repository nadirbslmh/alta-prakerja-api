package controllers

import (
	"context"
	"gugcp/models"
	"gugcp/services"
	"gugcp/utils"
	"math"
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

func GetAllTasks(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "invalid page parameter",
		})
	}

	if page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "invalid limit parameter",
		})
	}

	if limit <= 0 {
		limit = 10
	}

	username := c.QueryParam("username")

	ctx := context.WithValue(c.Request().Context(), utils.PageKey, page)
	ctx = context.WithValue(ctx, utils.LimitKey, limit)
	ctx = context.WithValue(ctx, utils.UsernameKey, username)

	tasks, err := services.GetAllTasks(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	totalTasks, err := services.CountTasks(c.Request().Context(), username)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: "error when calculating tasks",
		})
	}

	totalPages := int(math.Ceil(float64(totalTasks) / float64(limit)))

	responseData := struct {
		Status      bool              `json:"status"`
		Message     string            `json:"message"`
		Data        []models.TaskData `json:"data"`
		TotalPages  int               `json:"total_pages"`
		CurrentPage int               `json:"current_page"`
		NextPage    int               `json:"next_page,omitempty"`
		PrevPage    int               `json:"prev_page,omitempty"`
	}{
		Status:      true,
		Message:     "all tasks",
		Data:        tasks,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	if page < totalPages {
		responseData.NextPage = page + 1
	}

	if page > 1 {
		responseData.PrevPage = page - 1
	}

	return c.JSON(http.StatusOK, responseData)
}

func GetTaskByID(c echo.Context) error {
	taskID := c.Param("taskID")

	tID, err := strconv.Atoi(taskID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "task ID is invalid",
		})
	}

	task, err := services.RetrieveTaskByID(c.Request().Context(), tID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[models.Task]{
		Status:  true,
		Message: "task data found",
		Data:    task,
	})
}
