package services

import (
	"gugcp/utils"
	"mime/multipart"
)

func Upload(file *multipart.FileHeader) (string, error) {
	//TODO: upload to actual storage (GCS)
	filename, err := utils.UploadToStorage(file)

	if err != nil {
		return "", err
	}

	return filename, err
	//TODO: send file URL to Prakerja API
}
