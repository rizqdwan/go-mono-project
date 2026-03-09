package http

import (
	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/internal/auth"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)


func SetupRouter(e *echo.Echo, tokenService *token.Service, authMiddleware *auth.AuthMiddleware) {

	api := e.Group("/api/v1")

	// authGroup := api.Group("/auth")
	// authGroup.POST("/signin", authHandler.Login)
	// authGroup.POST("/renew", authHandler.RenewToken)

	// protected := api.Group("", middleware.AuthMiddleware(tokenService))

	// authProtected := protected.Group("/auth")
	// authProtected.POST("/signup", authHandler.Register,
	// 	middleware.RoleMiddleware("PROJECT_ADMIN"))
	// authProtected.POST("/logout", authHandler.Logout)
	// authProtected.PUT("/change-password", authHandler.ChangePassword)
}