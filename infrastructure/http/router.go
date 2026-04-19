package http

import (
	"github.com/labstack/echo/v5"
	_ "github.com/rizqdwan/go-mono-project/docs"
	"github.com/rizqdwan/go-mono-project/infrastructure/http/middleware"
	"github.com/rizqdwan/go-mono-project/internal/auth"
	"github.com/rizqdwan/go-mono-project/internal/user/role"
	"github.com/rizqdwan/go-mono-project/pkg/token"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func SetupRouter(e *echo.Echo, tokenSvc *token.Service, authHandler *auth.Handler) {

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")

	welcomeHandler := NewWelcomeHandler()
	api.GET("", welcomeHandler.Welcome)

	registerAuthRoutes(api, tokenSvc, authHandler)

	// registerUserRoutes(api, tokenSvc, userHandler)
	// registerDepartmentRoutes(api, tokenSvc, DepartmentHandler)
	// registerGroupRoutes(api, tokenSvc, GroupHandler)
}

func registerAuthRoutes(api *echo.Group, tokenSvc *token.Service, h *auth.Handler) {
	authGroup := api.Group("/auth")

	authGroup.POST("/signin", h.Login)
	authGroup.POST("/renew", h.RenewToken)
	authGroup.POST("/logout", h.Logout)

	protected := authGroup.Group("", middleware.AuthMiddleware(tokenSvc))
	protected.PUT("/change-password", h.ChangePassword)
	protected.POST("/signup", h.Register, middleware.RolesMiddleware(role.ProjectAdmin))
}
