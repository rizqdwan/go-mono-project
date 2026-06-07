package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
)

const RequestIDKey = "request_id"

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestID := c.Request().Header.Get("X-Request-ID")

			if requestID == "" {
				requestID = generateID()
			}

			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)

			return next(c)
		}
	}
}

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			method := c.Request().Method
			path := c.Request().URL.Path

			err := next(c)
			status := http.StatusOK

			if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
				status = resp.Status
			}

			if err != nil {
				var appErr commonErrors.AppError
				if errors.As(err, &appErr) {
					status = appErr.HttpCode
				} else if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else {
					status = http.StatusInternalServerError
				}
			}

			requestID, _ := c.Get(RequestIDKey).(string)

			slog.Info("request",
				"request_id", requestID,
				"method", method,
				"path", path,
				"status", status,
				"latency_ms", time.Since(start).Milliseconds(),
			)
			return err
		}
	}
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
