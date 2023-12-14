package services

import (
	"bytes"
	"context"
	"database/sql"
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

func SaveRedeemCode(ctx context.Context, input models.RedeemInput) (models.Redeem, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when creating transaction: %v", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO wpone_prakerja_redeems(user_id,state,redeem_code,sequence,status) VALUES (?,?,?,?,?)",
		input.UserID, input.State, input.RedeemCode, input.Sequence, 0,
	)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when saving redeem code: %v", err)
	}

	var redeem models.Redeem

	result := tx.QueryRowContext(ctx, "SELECT * FROM wpone_prakerja_redeems WHERE state = ?", input.State)

	if err := result.Scan(&redeem.ID, &redeem.UserID, &redeem.State, &redeem.RedeemCode, &redeem.Sequence, &redeem.Status); err != nil {
		if err == sql.ErrNoRows {
			return models.Redeem{}, fmt.Errorf("redeem is not exists: %v", err)
		}
		return models.Redeem{}, fmt.Errorf("error when getting redeem: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Redeem{}, fmt.Errorf("error when starting transaction: %v", err)
	}

	return redeem, nil
}

func CheckAttendanceStatus(ctx context.Context, input models.CheckStatusInput) (models.Redeem, error) {
	checkResult, err := getAttendanceStatus(input)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when checking status: %v", err)
	}

	if checkResult.Data.AttendanceStatus != 1 {
		return models.Redeem{}, fmt.Errorf("attendance status is invalid")
	}

	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when creating transaction: %v", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"UPDATE wpone_prakerja_redeems SET status=? WHERE state=?",
		1, input.State,
	)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when updating redeem status: %v", err)
	}

	var redeem models.Redeem

	result := tx.QueryRowContext(ctx, "SELECT * FROM wpone_prakerja_redeems WHERE state = ?", input.State)

	if err := result.Scan(&redeem.ID, &redeem.UserID, &redeem.State, &redeem.RedeemCode, &redeem.Sequence, &redeem.Status); err != nil {
		if err == sql.ErrNoRows {
			return models.Redeem{}, fmt.Errorf("redeem is not exists: %v", err)
		}
		return models.Redeem{}, fmt.Errorf("error when getting redeem: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Redeem{}, fmt.Errorf("error when starting transaction: %v", err)
	}

	return redeem, nil
}

func getAttendanceStatus(input models.CheckStatusInput) (models.CheckStatusResponse, error) {
	url := "https://api.prakerja.go.id/api/v1/integration/payment/redeem-code/status"

	clientCode := utils.GetConfig("CLIENT_CODE")
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))
	endpoint := "/api/v1/integration/payment/redeem-code/status"
	method := http.MethodPost

	reqInput := models.CheckStatusRequest{
		RedeemCode: input.RedeemCode,
		Sequence:   input.Sequence,
	}

	signature, err := utils.GenerateSignature(reqInput, timestamp, endpoint, method)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
		return models.CheckStatusResponse{}, err
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","sequence":%d}`, input.RedeemCode, input.Sequence))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
		return models.CheckStatusResponse{}, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("client_code", clientCode)
	req.Header.Set("signature", signature)
	req.Header.Set("timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		log.Printf("error when sending HTTP request: %v", err)
		return models.CheckStatusResponse{}, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("error when parsing request body: %v", err)
		return models.CheckStatusResponse{}, err
	}

	response, err := models.UnmarshalCheckStatusResponse(body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.CheckStatusResponse{}, err
	}

	return response, nil
}
