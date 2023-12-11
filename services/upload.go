package services

import (
	"gugcp/models"
	"gugcp/utils"
)

func Upload(uploadDTO models.UploadDTO) (string, error) {
	//TODO: upload to actual storage (GCS)
	filename, err := utils.UploadToStorage(uploadDTO.File)

	if err != nil {
		return "", err
	}

	return filename, err
	//TODO: send file URL to Prakerja API
}
