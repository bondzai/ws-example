package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalStorage implements the FileStorage interface for saving files to the local disk.
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new LocalStorage instance.
// - basePath is the directory where files will be stored (e.g., "./uploads").
// - baseURL is the public URL prefix for accessing the files (e.g., "http://localhost:8080/files").
func NewLocalStorage(basePath, baseURL string) (FileStorage, error) {
	// Ensure the base path exists.
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Upload saves a file to the local disk and returns its public URL.
func (s *LocalStorage) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	// Generate a unique filename to prevent collisions.
	uniqueFileName := uuid.New().String() + filepath.Ext(fileName)
	filePath := filepath.Join(s.basePath, uniqueFileName)

	// Create the file.
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the content from the reader to the file.
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	// Return the public URL.
	url := fmt.Sprintf("%s/%s", s.baseURL, uniqueFileName)
	return url, nil
}
