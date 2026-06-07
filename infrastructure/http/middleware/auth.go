package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)

func AuthMiddleware(tokenService *token.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authorization header")
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenService.ValidateAccessToken(tokenStr)
			if err != nil {
				switch err {
				case token.ErrExpiredToken:
					return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
				default:
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
				}
			}

			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RolesMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role == "" {
				return echo.NewHTTPError(http.StatusForbidden, "no role found in context")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

func GetUserID(c *echo.Context) (int64, error) {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing user identity")
	}
	return userID, nil
}
