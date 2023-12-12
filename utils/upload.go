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

const (
	baseURL    = "https://storage.googleapis.com/"
	projectID  = "alta-prakerja"
	bucketName = "alta-prakerja"
	uploadPath = "tasks/"
)

var client *storage.Client

func UploadToStorage(file *multipart.FileHeader) (string, error) {
	var err error

	client, err = storage.NewClient(context.Background(), option.WithCredentialsFile("accesskey.json"))

	if err != nil {
		return "", errors.New("error when creating storage client")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).IAM().V3().Policy: %w", bucketName, err)
	}
	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})
	if err := client.Bucket(bucketName).IAM().V3().SetPolicy(ctx, policy); err != nil {
		return "", fmt.Errorf("Bucket(%q).IAM().SetPolicy: %w", bucketName, err)
	}

	blobFile, err := file.Open()

	if err != nil {
		return "", errors.New("open file failed")
	}

	uploadedFileName := GenerateFileName(file.Filename)

	objectName := uploadPath + uploadedFileName

	sw := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(sw, blobFile); err != nil {
		return "", errors.New("error when copying file")
	}

	if err := sw.Close(); err != nil {
		return "", errors.New("error when closing storage client")
	}

	publicURL := getPublicURL(objectName)

	log.Println("upload complete! ", publicURL)

	return publicURL, nil
}

func CloseStorageClient() error {
	return client.Close()
}

func getPublicURL(objectName string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, projectID, objectName)
}
