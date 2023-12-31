package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func getAllowedExtensions() []string {
	return []string{".pdf", ".docx"}
}

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

func ValidateFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return slices.Contains(getAllowedExtensions(), ext)
}

func GetConfig(key string) string {
	if os.Getenv("MODE") == "production" {
		return os.Getenv(key)
	}

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error when reading configuration file: %s\n", err)
	}

	return viper.GetString(key)
}

func GenerateSignature(input any, timestamp int64, endpoint, method string) (string, error) {
	clientCode := GetConfig("CLIENT_CODE")
	signKey := GetConfig("SIGNATURE_KEY")

	jsonStr, err := json.Marshal(input)

	if err != nil {
		return "", err
	}

	requestInput := clientCode + fmt.Sprintf("%d", timestamp) + method + endpoint + string(jsonStr)

	h := hmac.New(sha1.New, []byte(signKey))
	h.Write([]byte(requestInput))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature, nil
}
