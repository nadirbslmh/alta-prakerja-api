package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

// TODO: upload to actual storage (GCS)
func UploadToStorage(file *multipart.FileHeader) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", errors.New("could not open the file")
	}

	defer src.Close()

	// Create a unique filename for the uploaded file
	uploadedFileName := GenerateFileName(file.Filename)

	// Create the destination file
	dst, err := os.Create("uploads/" + uploadedFileName)
	if err != nil {
		return "", errors.New("could not create the destination file")
	}
	defer dst.Close()

	// Copy the contents of the source file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return "", errors.New("could not copy file content")
	}

	return uploadedFileName, nil
}
