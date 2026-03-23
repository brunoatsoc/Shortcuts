package httpserver

import (
	"intellisys-backend/internal/api/controller"
	"intellisys-backend/internal/api/middleware"
	"intellisys-backend/internal/domain/user"
	"intellisys-backend/internal/infra/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	shortcutController *controller.ShortcutController,
	userConstroller *controller.UserController,
	authController *controller.AuthController,
	authMiddleware *middleware.AuthMidware,
	cfg *config.Config,
) *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", cfg.FalePath.Path)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RateLimitMiddleware())

	setupShortcutRoutes(router, shortcutController, authMiddleware)
	setupUsuarioRoutes(router, userConstroller, authMiddleware)
	setupAuthRoutes(router, authController)

	return router
}

func setupShortcutRoutes(
	router *gin.Engine,
	shortcutController *controller.ShortcutController,
	authMiddleware *middleware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		shortcutRoutes := v1.Group("/shortcuts")
		{
			authenticatedRoutes := v1.Group("/admin")
			authenticatedRoutes.Use(authMiddleware.AuthenticationMidware(), middleware.RoleMidware(user.RoleAdmin))
			{
				authenticatedRoutes.POST("", shortcutController.Create)
				authenticatedRoutes.PUT("/:id", shortcutController.Update)
				authenticatedRoutes.DELETE("/:id", shortcutController.Delete)
			}

			shortcutRoutes.GET("", shortcutController.GetAll)
		}
	}
}

func setupUsuarioRoutes(
	router *gin.Engine,
	userController *controller.UserController,
	authMidware *middleware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(authMidware.AuthenticationMidware(), middleware.RoleMidware(user.RoleAdmin))
		{
			adminRoutes.GET("/all", userController.GetAll)
			adminRoutes.GET("", userController.GetUserById)
			adminRoutes.PUT("", userController.UpdateUsuario)
			adminRoutes.DELETE("", userController.DeleteUser)
		}
	}
}

func setupAuthRoutes(
	router *gin.Engine,
	authController *controller.AuthController,
) {
	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", authController.Login)
		}
	}
}
