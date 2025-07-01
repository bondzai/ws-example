package errs

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomError struct {
	Message string
	Code    int
}

func (e CustomError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return CustomError{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func NewNotFoundError(message string) error {
	return CustomError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewInternalServerError(message string) error {
	return CustomError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewUnexpectedError() error {
	return CustomError{
		Message: "An unexpected error occurred",
		Code:    http.StatusInternalServerError,
	}
}

func HandleError(err error) error {
	switch e := err.(type) {
	case CustomError:
		return &fiber.Error{
			Code:    e.Code,
			Message: e.Message,
		}
	default:
		if err == gorm.ErrRecordNotFound {
			return NewNotFoundError(err.Error())
		}

		if strings.Contains(err.Error(), "SQLSTATE") {
			return NewInternalServerError(err.Error())
		}

		return NewUnexpectedError()
	}
}

func HandleFiberError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case CustomError:
		return c.Status(e.Code).JSON(fiber.Map{
			"message": e.Message,
		})
	default:
		return NewUnexpectedError()
	}
}
