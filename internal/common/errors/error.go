package errors

import "net/http"

type AppError struct {
	Message  string
	HttpCode int
	Err      error
}

func (e AppError) Error() string {
	return e.Message
}
func (e AppError) Unwrap() error { return e.Err }

func InternalServerError(message string, err error) AppError {
	return AppError{
		Message:  message,
		HttpCode: http.StatusInternalServerError,
		Err:      err,
	}
}

func InvariantError(message string, err error) AppError {
	return AppError{
		Message:  message,
		HttpCode: http.StatusBadRequest,
		Err:      err,
	}
}

func NotFoundError(message string, err error) AppError {
	return AppError{
		Message:  message,
		HttpCode: http.StatusNotFound,
		Err:      err,
	}
}

func UnauthorizedError(message string, err error) AppError {
	return AppError{
		Message:  message,
		HttpCode: http.StatusUnauthorized,
		Err:      err,
	}
}
