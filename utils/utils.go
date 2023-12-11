package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error when reading configuration file: %s\n", err)
	}

	return viper.GetString(key)
}

func GenerateSignature(input any, timestamp int64, endpoint, method string) (string, error) {
	clientCode := "alterra-academy"
	signKey := "db6a42a727104cd6a887b73df599ea29"

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
