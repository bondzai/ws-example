package filestorage

import (
	"context"
	"io"
)

// FileStorage defines the interface for a file storage service.
// This abstraction allows for easy swapping between local and cloud-based storage.
type FileStorage interface {
	// Upload saves a file from an io.Reader and returns its public-facing URL.
	Upload(ctx context.Context, reader io.Reader, fileName string) (string, error)
}
