package api

import (
	"database/sql"
	"intellisys-backend/internal/configuration"
	"intellisys-backend/internal/controllers"
	"intellisys-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ServerAPI struct {
	db *sql.DB
}

func NewServerAPI(db *sql.DB) *ServerAPI {
	return &ServerAPI{db: db}
}

func (s *ServerAPI) RunServer() error {
	r := gin.Default()

	/*r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // URL do Vite React
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))*/

	shortcutService := service.NewShortcutService(s.db)
	shortcutController := controllers.NewShortcutController(shortcutService)
	shortcutController.RegisterRoutes(r)

	userService := service.NewUsersService(s.db)
	userController := controllers.NewUsersController(userService)
	userController.RegisterRoutes(r)

	return r.Run(configuration.Config.ServerPort)
}
