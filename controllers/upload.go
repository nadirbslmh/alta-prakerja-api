package controllers

import (
	"gugcp/models"
	"gugcp/services"
	"gugcp/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
	var uploadForm models.UploadForm

	if err := c.Bind(&uploadForm); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "invalid request",
		})
	}

	if err := uploadForm.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response[[]*models.ValidationErrorResponse]{
			Status:  false,
			Message: "request validation failed",
			Data:    err,
		})
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "missing file parameter",
		})
	}

	filename := file.Filename

	// Validate the file
	isValidFile := utils.ValidateFile(filename)

	if !isValidFile {
		return c.JSON(http.StatusBadRequest, models.Response[any]{
			Status:  false,
			Message: "file type is invalid",
		})
	}

	uploadDTO := models.UploadDTO{
		File:           file,
		UploadFormData: uploadForm,
	}

	res, err := services.Upload(c.Request().Context(), uploadDTO)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[models.UploadResponse]{
		Status:  true,
		Message: "file uploaded successfully",
		Data:    res,
	})
}
