package utils

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
)

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
