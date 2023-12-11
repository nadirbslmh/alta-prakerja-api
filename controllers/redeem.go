package controllers

import (
	"gugcp/models"
	"gugcp/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SaveRedeemCode(c echo.Context) error {
	var input models.RedeemInput

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

	redeem, err := services.SaveRedeemCode(c.Request().Context(), input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, models.Response[models.Redeem]{
		Status:  true,
		Message: "redeem saved successfully",
		Data:    redeem,
	})
}

func CheckStatus(c echo.Context) error {
	var input models.CheckStatusInput

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

	redeem, err := services.CheckAttendanceStatus(c.Request().Context(), input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[models.Redeem]{
		Status:  true,
		Message: "redeem status updated successfully",
		Data:    redeem,
	})
}
