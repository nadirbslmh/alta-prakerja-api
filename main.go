package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Define a route for file upload
	e.POST("/upload", uploadFile)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler function for file upload
func uploadFile(c echo.Context) error {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing file parameter",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not open the file",
		})
	}
	defer src.Close()

	// Create a unique filename for the uploaded file
	uploadedFileName := generateFileName(file.Filename)

	// Create the destination file
	dst, err := os.Create("uploads/" + uploadedFileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not create the destination file",
		})
	}
	defer dst.Close()

	// Copy the contents of the source file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not copy file content",
		})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "File uploaded successfully",
		"filename": uploadedFileName,
	})
}

// Function to generate a unique filename
func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
