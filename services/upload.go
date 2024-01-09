package services

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gugcp/database"
	"gugcp/models"
	"gugcp/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Upload(ctx context.Context, uploadDTO models.UploadDTO) (models.UploadResponse, error) {
	fileURL, err := utils.UploadToStorage(uploadDTO.File)

	if err != nil {
		return models.UploadResponse{}, err
	}

	redeem, err := getRedeemByUserID(ctx, uploadDTO.UploadFormData.UserID)

	if err != nil {
		return models.UploadResponse{}, errors.New("redeem not found")
	}

	userRedeemCode := redeem.RedeemCode

	request := models.UploadRequest{
		RedeemCode: userRedeemCode,
		Scope:      uploadDTO.UploadFormData.Scope,
		Sequence:   uploadDTO.UploadFormData.Sequence,
		FileURL:    fileURL,
	}

	res, err := submitTask(request)

	if err != nil {
		return models.UploadResponse{}, err
	}

	uploadDTO.UploadFormData.RedeemCode = userRedeemCode

	err = saveTaskToDB(ctx, uploadDTO, fileURL)

	if err != nil {
		return models.UploadResponse{}, err
	}

	return res, nil
}

func submitTask(request models.UploadRequest) (models.UploadResponse, error) {
	url := "https://api.prakerja.go.id/api/v1/integration/tpm/submission"

	clientCode := utils.GetConfig("CLIENT_CODE")
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))
	endpoint := "/api/v1/integration/tpm/submission"
	method := http.MethodPost

	signature, err := utils.GenerateSignature(request, timestamp, endpoint, method)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
		return models.UploadResponse{}, errors.New("error when creating signature")
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","scope":"%s","sequence":%d,"url_file":"%s"}`,
		request.RedeemCode,
		request.Scope,
		request.Sequence,
		request.FileURL,
	))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
		return models.UploadResponse{}, errors.New("error when creating HTTP request")
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("client_code", clientCode)
	req.Header.Set("signature", signature)
	req.Header.Set("timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		log.Printf("error when sending HTTP request: %v", err)
		return models.UploadResponse{}, errors.New("error when sending HTTP request")
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.UploadResponse{}, errors.New("error when parsing response")
	}

	response, err := models.UnmarshalUploadResponse(body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.UploadResponse{}, errors.New("error when parsing response")
	}

	if result.StatusCode != http.StatusOK {
		log.Printf("error when submitting task: %v", response.Message)
		return models.UploadResponse{}, errors.New(response.Message)
	}

	return response, nil
}

func saveTaskToDB(ctx context.Context, uploadDTO models.UploadDTO, fileURL string) error {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO wpone_prakerja_task(user_ID,sesi,sequence,link,batch,redeem_code,scope) VALUES (?,?,?,?,?,?,?)",
		uploadDTO.UploadFormData.UserID,
		uploadDTO.UploadFormData.Session,
		uploadDTO.UploadFormData.Sequence,
		fileURL,
		uploadDTO.UploadFormData.Batch,
		uploadDTO.UploadFormData.RedeemCode,
		uploadDTO.UploadFormData.Scope,
	)

	if err != nil {
		log.Printf("error when saving task data: %v", err)
		return errors.New("error when saving task data")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return errors.New("error when starting transaction")
	}

	return nil
}

func getRedeemByUserID(ctx context.Context, userID int) (models.Redeem, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return models.Redeem{}, errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	var redeem models.Redeem

	result := tx.QueryRowContext(ctx, "SELECT * FROM wpone_prakerja_redeems WHERE user_id = ?", userID)

	if err := result.Scan(&redeem.ID, &redeem.UserID, &redeem.State, &redeem.RedeemCode, &redeem.Sequence, &redeem.Status); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("redeem is not exists: %v", err)
			return models.Redeem{}, errors.New("redeem is not exists")
		}
		log.Printf("error when getting redeem: %v", err)
		return models.Redeem{}, errors.New("error when getting redeem")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return models.Redeem{}, errors.New("error when starting transaction")
	}

	return redeem, nil
}
