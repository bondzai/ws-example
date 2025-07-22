package usecases

import (
	"api-gateway/pkg/filestorage"
	"context"
	"io"
)

// FileUploadResponse is the DTO returned after a successful file upload.
type FileUploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}

// FileUploadUseCase defines the business logic for uploading files.
type FileUploadUseCase interface {
	// UploadFile handles the business logic of storing a file and returning its metadata.
	UploadFile(ctx context.Context, reader io.Reader, fileName string, fileSize int64) (*FileUploadResponse, error)
}

// fileUploadUseCase implements the FileUploadUseCase.
type fileUploadUseCase struct {
	fileStorage filestorage.FileStorage
}

// NewFileUploadUseCase creates a new FileUploadUseCase.
func NewFileUploadUseCase(fileStorage filestorage.FileStorage) FileUploadUseCase {
	return &fileUploadUseCase{
		fileStorage: fileStorage,
	}
}

// UploadFile saves the file using the configured file storage and returns the URL.
func (uc *fileUploadUseCase) UploadFile(ctx context.Context, reader io.Reader, fileName string, fileSize int64) (*FileUploadResponse, error) {
	url, err := uc.fileStorage.Upload(ctx, reader, fileName)
	if err != nil {
		return nil, err
	}

	return &FileUploadResponse{
		URL:      url,
		FileName: fileName,
		FileSize: fileSize,
	}, nil
}
