package controllers

import (
	"gugcp/models"
	"gugcp/services"
	"gugcp/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
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

	res, err := services.Upload(file)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response[any]{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response[string]{
		Status:  true,
		Message: "file uploaded successfully",
		Data:    res,
	})

	// // Open the file
	// src, err := file.Open()
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, models.Response[any]{
	// 		Status:  false,
	// 		Message: "could not open the file",
	// 	})
	// }
	// defer src.Close()

	// // Create a unique filename for the uploaded file
	// uploadedFileName := utils.GenerateFileName(file.Filename)

	// // Create the destination file
	// dst, err := os.Create("uploads/" + uploadedFileName)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, models.Response[any]{
	// 		Status:  false,
	// 		Message: "could not create the destination file",
	// 	})
	// }
	// defer dst.Close()

	// // Copy the contents of the source file to the destination file
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, models.Response[any]{
	// 		Status:  false,
	// 		Message: "could not copy file content",
	// 	})
	// }
}
