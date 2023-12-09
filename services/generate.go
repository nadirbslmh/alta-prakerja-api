package services

import (
	"bytes"
	"fmt"
	"gugcp/models"
	"gugcp/utils"
	"io"
	"net/http"
	"strconv"
	"time"

	"log"
)

func GenerateURL(input models.GenerateInput) (models.GenerateURLResponse, error) {
	url := "https://api.prakerja.go.id/api/v1/integration/oauth/url-generate"

	clientCode := "alterra-academy"
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))
	endpoint := "/api/v1/integration/oauth/url-generate"
	method := http.MethodPost

	signature, err := utils.GenerateSignature(input, timestamp, endpoint, method)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
		return models.GenerateURLResponse{}, err
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","sequence":%d,"redirect_uri":"%s","email":"%s"}`, input.RedeemCode, input.Sequence, input.RedirectURI, input.Email))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
		return models.GenerateURLResponse{}, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Client_code", clientCode)
	req.Header.Set("Signature", signature)
	req.Header.Set("Timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		log.Printf("error when sending HTTP request: %v", err)
		return models.GenerateURLResponse{}, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("error when parsing request body: %v", err)
		return models.GenerateURLResponse{}, err
	}

	response, err := models.UnmarshalGenerateURLResponse(body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.GenerateURLResponse{}, err
	}

	return response, nil
}
