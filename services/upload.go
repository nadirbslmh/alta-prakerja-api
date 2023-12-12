package services

import (
	"bytes"
	"fmt"
	"gugcp/models"
	"gugcp/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Upload(uploadDTO models.UploadDTO) (models.UploadResponse, error) {
	fileURL, err := utils.UploadToStorage(uploadDTO.File)

	if err != nil {
		return models.UploadResponse{}, err
	}

	request := models.UploadRequest{
		RedeemCode: uploadDTO.UploadRequestForm.RedeemCode,
		Scope:      uploadDTO.UploadRequestForm.Scope,
		Sequence:   uploadDTO.UploadRequestForm.Sequence,
		FileURL:    fileURL,
	}

	res, err := submitTask(request)

	if err != nil {
		return models.UploadResponse{}, err
	}

	return res, nil
}

func submitTask(request models.UploadRequest) (models.UploadResponse, error) {
	url := "https://api.prakerja.go.id/api/v1/integration/tpm/submission"

	clientCode := "alterra-academy"
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))
	endpoint := "/api/v1/integration/tpm/submission"
	method := http.MethodPost

	signature, err := utils.GenerateSignature(request, timestamp, endpoint, method)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
		return models.UploadResponse{}, err
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","scope":"%s","sequence":%d,"url_file":"%s"}`,
		request.RedeemCode,
		request.Scope,
		request.Sequence,
		request.FileURL,
	))

	log.Println("request data: ", string(data))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
		return models.UploadResponse{}, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("client_code", clientCode)
	req.Header.Set("signature", signature)
	req.Header.Set("timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		log.Printf("error when sending HTTP request: %v", err)
		return models.UploadResponse{}, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("error when parsing request body: %v", err)
		return models.UploadResponse{}, err
	}

	response, err := models.UnmarshalUploadResponse(body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.UploadResponse{}, err
	}

	return response, nil
}
