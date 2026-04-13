package httpserver

import (
	"intellisys-backend/internal/api/controller"
	"intellisys-backend/internal/api/middleware"
	"intellisys-backend/internal/domain/log"
	"intellisys-backend/internal/domain/user"
	"intellisys-backend/internal/infra/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In

	ShortcutController *controller.ShortcutController
	UserController     *controller.UserController
	AuthController     *controller.AuthController
	AuthMiddleware     *middleware.AuthMidware
	Cfg                *config.Config
	Log                log.Logger
	GlobalLimiter      *middleware.RateLimiter `name:"globalLimiter"`
	LoginLimiter       *middleware.RateLimiter `name:"loginLimiter"`
}

func NewRouter(p RouterParams) *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", p.Cfg.FalePath.Path)

	// Para produção com APACHE essa configuração de CORs deve ser comentada
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(p.GlobalLimiter.Middleware())
	router.Use(middleware.LoggingMiddleware(p.Log))

	setupShortcutRoutes(router, p.ShortcutController, p.AuthMiddleware)
	setupUsuarioRoutes(router, p.UserController, p.AuthMiddleware)
	setupAuthRoutes(router, p.AuthController, p.LoginLimiter)

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
	loginLimiter *middleware.RateLimiter,
) {
	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", loginLimiter.Middleware(), authController.Login)
		}
	}
}
