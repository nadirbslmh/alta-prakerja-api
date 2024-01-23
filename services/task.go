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

	request := models.FeedbackRequest{
		RedeemCode: task.RedeemCode,
		Scope:      task.Scope,
		Sequence:   int64(task.Sequence),
		Notes:      input.Notes,
		URLFile:    task.Link,
	}

	result, err := submitFeedback(request)

	if err != nil {
		return models.FeedbackResponse{}, err
	}

	err = updateTask(ctx, input, taskID)

	if err != nil {
		return models.FeedbackResponse{}, err
	}

	return result, nil
}

func RetrieveTaskByID(ctx context.Context, taskID int) (models.Task, error) {
	task, err := getTaskByID(ctx, taskID)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func GetAllTasks(ctx context.Context) ([]models.TaskData, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return []models.TaskData{}, errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	tasks := []models.TaskData{}
	task := models.TaskData{}

	page := ctx.Value(utils.PageKey).(int)
	limit := ctx.Value(utils.LimitKey).(int)
	username := ctx.Value(utils.UsernameKey).(string)
	batch := ctx.Value(utils.BatchKey).(string)

	offset := (page - 1) * limit

	var query string

	if username != "" {
		query = fmt.Sprintf(`SELECT wpone_prakerja_task.ID, wpone_prakerja_task.user_ID, wpone_prakerja_task.sequence, wpone_prakerja_task.link, wpone_prakerja_task.scope, wpone_prakerja_task.batch, wpone_prakerja_task.feedback, wpone_users.display_name FROM wpone_prakerja_task
            JOIN wpone_users ON wpone_prakerja_task.user_ID = wpone_users.ID 
            WHERE wpone_users.display_name LIKE '%%%s%%' AND batch = '%s' LIMIT %d OFFSET %d;`, username, batch, limit, offset)
	} else {
		query = fmt.Sprintf(`SELECT wpone_prakerja_task.ID, wpone_prakerja_task.user_ID, wpone_prakerja_task.sequence, wpone_prakerja_task.link, wpone_prakerja_task.scope, wpone_prakerja_task.batch, wpone_prakerja_task.feedback, wpone_users.display_name FROM wpone_prakerja_task
        JOIN wpone_users ON wpone_prakerja_task.user_ID = wpone_users.ID WHERE batch = '%s' LIMIT %d OFFSET %d;`, batch, limit, offset)
	}

	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		log.Printf("error when fetching tasks: %v", err)
		return []models.TaskData{}, errors.New("error when fetching tasks")
	}

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.UserID, &task.Sequence, &task.Link, &task.Scope, &task.Batch, &task.Feedback, &task.Name)
		if err != nil {
			log.Printf("error when fetching tasks: %v", err)
			return []models.TaskData{}, errors.New("error when fetching tasks")
		}
		tasks = append(tasks, task)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return []models.TaskData{}, errors.New("error when starting transaction")
	}

	return tasks, nil
}

func CountTasks(ctx context.Context) (int, error) {
	var query string

	username := ctx.Value(utils.UsernameKey).(string)
	batch := ctx.Value(utils.BatchKey).(string)

	if username != "" {
		query = fmt.Sprintf(`SELECT COUNT(*) FROM wpone_prakerja_task JOIN wpone_users ON wpone_prakerja_task.user_ID = wpone_users.ID WHERE wpone_users.display_name LIKE '%%%s%%' AND batch = '%s'`, username, batch)
	} else {
		query = fmt.Sprintf(`SELECT COUNT(*) FROM wpone_prakerja_task JOIN wpone_users ON wpone_prakerja_task.user_ID = wpone_users.ID WHERE batch = '%s'`, batch)
	}

	rows, err := database.DB.QueryContext(ctx, query)

	var count int

	if err != nil {
		return -1, err
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func GetTaskDetails(ctx context.Context, taskID int) (models.TaskData, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return models.TaskData{}, errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	var task models.TaskData

	result := tx.QueryRowContext(ctx, `SELECT wpone_prakerja_task.ID, wpone_prakerja_task.user_ID, wpone_prakerja_task.sequence, wpone_prakerja_task.link, wpone_prakerja_task.scope, wpone_prakerja_task.batch, wpone_prakerja_task.feedback, wpone_users.display_name FROM wpone_prakerja_task 
	JOIN wpone_users ON wpone_prakerja_task.user_ID = wpone_users.ID  WHERE wpone_prakerja_task.ID = ?`, taskID)

	if err := result.Scan(&task.ID, &task.UserID, &task.Sequence, &task.Link, &task.Scope, &task.Batch, &task.Feedback, &task.Name); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("task is not exists: %v", err)
			return models.TaskData{}, errors.New("task is not exists")
		}
		log.Printf("error when getting task: %v", err)
		return models.TaskData{}, errors.New("error when getting task")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return models.TaskData{}, errors.New("error when starting transaction")
	}

	return task, nil
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

	if err := result.Scan(&task.ID, &task.UserID, &task.Session, &task.Sequence, &task.Link, &task.Batch, &task.RedeemCode, &task.Scope, &task.Feedback); err != nil {
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

func updateTask(ctx context.Context, input models.FeedbackInput, taskID int) error {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("error when creating transaction: %v", err)
		return errors.New("error when creating transaction")
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"UPDATE wpone_prakerja_task SET feedback=? WHERE ID=?",
		input.Notes, taskID,
	)

	if err != nil {
		log.Printf("error when updating task: %v", err)
		return errors.New("error when updating task")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error when starting transaction: %v", err)
		return errors.New("error when starting transaction")
	}

	return nil
}
