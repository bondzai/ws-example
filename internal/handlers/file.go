package handlers

import (
	"api-gateway/internal/usecases"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// FileUploadHandler handles HTTP requests for file uploads.
type FileUploadHandler struct {
	useCase usecases.FileUploadUseCase
}

// NewFileUploadHandler creates a new FileUploadHandler.
func NewFileUploadHandler(useCase usecases.FileUploadUseCase) *FileUploadHandler {
	return &FileUploadHandler{
		useCase: useCase,
	}
}

// UploadFile is the handler for the POST /files/upload endpoint.
func (h *FileUploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	// Open the file for reading.
	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not open file"})
	}
	defer src.Close()

	// Upload the file using the use case.
	uploadRes, err := h.useCase.UploadFile(c.Context(), src, file.Filename, file.Size)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(uploadRes)
}
