package utils

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
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
