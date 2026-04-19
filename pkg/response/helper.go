package response

import (
	"github.com/labstack/echo/v5"
)

type WebResponse[T any] struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`
	Data     T              `json:"data,omitempty"`
	Errors   []string       `json:"errors,omitempty"`
	MetaData map[string]any `json:"metaData,omitempty"`
}

func Success[T any](c *echo.Context, status int, message string, data T) error {
	return c.JSON(status, WebResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func SuccessWithData[T any](c *echo.Context, status int, message string, data T, meta map[string]any) error {
	return c.JSON(status, WebResponse[T]{
		Status:   status,
		Message:  message,
		Data:     data,
		MetaData: meta,
	})
}

func Error(c *echo.Context, status int, message string, errors []string) error {
	return c.JSON(status, WebResponse[any]{
		Status:  status,
		Message: message,
		Errors:  errors,
	})
}

func PaginationMeta(page, size, totalElements, totalPages int, last bool) map[string]any {
	return map[string]any{
		"page":          page,
		"size":          size,
		"totalElements": totalElements,
		"totalPages":    totalPages,
		"last":          last,
	}
}
