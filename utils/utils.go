package utils

import (
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
