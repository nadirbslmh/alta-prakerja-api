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

func GenerateURL(input models.GenerateInput) {
	url := "https://api.prakerja.go.id/api/v1/integration/oauth/url-generate"

	clientCode := "alterra-academy"
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))

	signature, err := utils.GenerateSignature(input, timestamp)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","sequence":%d,"redirect_uri":"%s","email":"%s"}`, input.RedeemCode, input.Sequence, input.RedirectURI, input.Email))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Client_code", clientCode)
	req.Header.Set("Signature", signature)
	req.Header.Set("Timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer result.Body.Close()

	body, _ := io.ReadAll(result.Body)
	fmt.Println("response:", string(body))
}
