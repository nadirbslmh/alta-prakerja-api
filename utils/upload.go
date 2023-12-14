package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/iam"
	iampb "cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	baseURL    = GetConfig("STORAGE_BASE_URL")
	projectID  = GetConfig("STORAGE_PROJECT_ID")
	bucketName = GetConfig("STORAGE_BUCKET_NAME")
	uploadPath = GetConfig("STORAGE_UPLOAD_PATH")
)

var client *storage.Client

func InitStorageClient() {
	var err error

	client, err = storage.NewClient(context.Background(), option.WithCredentialsFile("accesskey.json"))

	if err != nil {
		log.Fatalf("error when creating storage client: %v", err)
	}

	log.Println("storage client initialized")
}

func UploadToStorage(file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		log.Printf("error when setting up a bucket: %v", err)
		return "", errors.New("error when setting up a bucket")
	}
	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})
	if err := client.Bucket(bucketName).IAM().V3().SetPolicy(ctx, policy); err != nil {
		log.Printf("error when setting up a bucket: %v", err)
		return "", errors.New("error when setting up a bucket")
	}

	blobFile, err := file.Open()

	if err != nil {
		log.Printf("open file failed: %v", err)
		return "", errors.New("open file failed")
	}

	uploadedFileName := GenerateFileName(file.Filename)

	objectName := uploadPath + uploadedFileName

	sw := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(sw, blobFile); err != nil {
		log.Printf("error when copying file: %v", err)
		return "", errors.New("error when copying file")
	}

	if err := sw.Close(); err != nil {
		log.Printf("error when closing storage client: %v", err)
		return "", errors.New("error when closing storage client")
	}

	publicURL := getPublicURL(objectName)

	return publicURL, nil
}

func CloseStorageClient() error {
	return client.Close()
}

func getPublicURL(objectName string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, projectID, objectName)
}
