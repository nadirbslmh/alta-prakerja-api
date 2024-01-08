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

func SendFeedback(ctx context.Context, input models.FeedbackInput, taskID int) (models.FeedbackResponse, error) {
	task, err := getTaskByID(ctx, taskID)

	if err != nil {
		return models.FeedbackResponse{}, err
	}

	seq, err := strconv.Atoi(task.Session)

	if err != nil {
		return models.FeedbackResponse{}, err
	}

	request := models.FeedbackRequest{
		RedeemCode: task.RedeemCode,
		Scope:      task.Scope,
		Sequence:   int64(seq),
		Notes:      input.Notes,
		URLFile:    task.Link,
	}

	result, err := submitFeedback(request)

	if err != nil {
		return models.FeedbackResponse{}, err
	}

	return result, nil
}

func getTaskByID(ctx context.Context, taskID int) (models.Task, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return models.Task{}, errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	var task models.Task

	result := tx.QueryRowContext(ctx, "SELECT * FROM wpone_prakerja_task WHERE ID = ?", taskID)

	if err := result.Scan(&task.ID, &task.UserID, &task.Session, &task.Link, &task.Batch, &task.RedeemCode, &task.Scope); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("task is not exists: %v", err)
			return models.Task{}, errors.New("task is not exists")
		}
		log.Printf("error when getting task: %v", err)
		return models.Task{}, errors.New("error when getting task")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return models.Task{}, errors.New("error when starting transaction")
	}

	return task, nil
}

func submitFeedback(request models.FeedbackRequest) (models.FeedbackResponse, error) {
	url := "https://api.prakerja.go.id/api/v1/integration/tpm/feedback/submission"

	clientCode := utils.GetConfig("CLIENT_CODE")
	contentType := "application/json"
	timestamp := time.Now().Unix()
	headerTimestamp := strconv.Itoa(int(timestamp))
	endpoint := "/api/v1/integration/tpm/feedback/submission"
	method := http.MethodPost

	signature, err := utils.GenerateSignature(request, timestamp, endpoint, method)

	if err != nil {
		log.Printf("error when creating signature: %v", err)
		return models.FeedbackResponse{}, errors.New("error when creating signature")
	}

	data := []byte(fmt.Sprintf(`{"redeem_code":"%s","scope":"%s","sequence":%d,"notes":"%s","url_file":"%s"}`,
		request.RedeemCode,
		request.Scope,
		request.Sequence,
		request.Notes,
		request.URLFile,
	))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error when creating HTTP request: %v", err)
		return models.FeedbackResponse{}, errors.New("error when creating HTTP request")
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("client_code", clientCode)
	req.Header.Set("signature", signature)
	req.Header.Set("timestamp", headerTimestamp)

	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		log.Printf("error when sending HTTP request: %v", err)
		return models.FeedbackResponse{}, errors.New("error when sending HTTP request")
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.FeedbackResponse{}, errors.New("error when parsing response")
	}

	response, err := models.UnmarshalFeedbackResponse(body)

	if err != nil {
		log.Printf("error when parsing response body: %v", err)
		return models.FeedbackResponse{}, errors.New("error when parsing response")
	}

	if result.StatusCode != http.StatusOK {
		log.Printf("error when submitting task: %v", response.Message)
		return models.FeedbackResponse{}, errors.New(response.Message)
	}

	return response, nil
}
