package http

import (
	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/config"
	_ "github.com/rizqdwan/go-mono-project/docs"
	"github.com/rizqdwan/go-mono-project/infrastructure/http/middleware"
	"github.com/rizqdwan/go-mono-project/internal/auth"
	"github.com/rizqdwan/go-mono-project/internal/organization/department"
	"github.com/rizqdwan/go-mono-project/internal/organization/group"
	"github.com/rizqdwan/go-mono-project/internal/user"
	"github.com/rizqdwan/go-mono-project/pkg/token"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

type Handlers struct {
	Auth       *auth.Handler
	User       *user.Handler
	Department *department.Handler
	Group      *group.Handler
}

func SetupRouter(e *echo.Echo, tokenSvc *token.Service, roles config.RoleConfig, h *Handlers) {

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")

	welcomeHandler := NewWelcomeHandler()
	api.GET("", welcomeHandler.Welcome)

	registerAuthRoutes(api, tokenSvc, roles, h.Auth)
	registerUserRoutes(api, tokenSvc, roles, h.User)
	registerOrganizationRoutes(api, tokenSvc, roles, h.Department, h.Group)
}

func registerAuthRoutes(api *echo.Group, tokenSvc *token.Service, roles config.RoleConfig, a *auth.Handler) {
	authGroup := api.Group("/auth")

	authGroup.POST("/signin", a.Login)
	authGroup.POST("/renew", a.RenewToken)
	authGroup.POST("/logout", a.Logout)

	protected := authGroup.Group("", middleware.AuthMiddleware(tokenSvc))
	protected.POST("/signup", a.Register, middleware.RolesMiddleware(roles.ProjectAdmin))
}

func registerUserRoutes(api *echo.Group, tokenSvc *token.Service, roles config.RoleConfig, u *user.Handler) {
	userGroup := api.Group("/users")

	protected := userGroup.Group("", middleware.AuthMiddleware(tokenSvc))
	protected.GET("/current", u.UserDetails)
	protected.PUT("/change-password", u.ChangePassword)

	admin := userGroup.Group("", middleware.AuthMiddleware(tokenSvc), middleware.RolesMiddleware(roles.ProjectAdmin))
	admin.GET("/list", u.UserList)
	admin.PUT("/:id/details", u.UpdateUserDetails)
	admin.PUT("/:id/reset-password", u.ResetPassword)
	admin.DELETE("/:id", u.DeleteUser)

}

func registerOrganizationRoutes(api *echo.Group, tokenSvc *token.Service, roles config.RoleConfig, d *department.Handler, g *group.Handler) {
	departments := api.Group("/departments", middleware.AuthMiddleware(tokenSvc), middleware.RolesMiddleware(roles.ProjectAdmin))

	departments.GET("", d.DepartmentList)
	departments.POST("", d.CreateDepartment)
	departments.PUT("/:id/details", d.UpdateDepartment)
	departments.DELETE("/:id", d.DeleteDepartment)

	groups := api.Group("/groups", middleware.AuthMiddleware(tokenSvc), middleware.RolesMiddleware(roles.ProjectAdmin))

	groups.GET("", g.GroupList)
	groups.POST("", g.CreateGroup)
	groups.PUT("/:id/details", g.UpdateGroup)
	groups.DELETE("/:id", g.DeleteGroup)

}
