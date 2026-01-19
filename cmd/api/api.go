package api

import (
	"database/sql"
	"intellisys-backend/internal/configuration"
	"intellisys-backend/internal/controllers"
	"intellisys-backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerAPI struct {
	db *sql.DB
}

func NewServerAPI(db *sql.DB) *ServerAPI {
	return &ServerAPI{db: db}
}

// Método para iniciar o servidor da API
func (s *ServerAPI) RunServer() error {
	r := gin.Default()

	// Configuração CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // URL do Vite React
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Registrar controladores e rotas para atalhos
	shortcutService := service.NewShortcutService(s.db)
	shortcutController := controllers.NewShortcutController(shortcutService)
	shortcutController.RegisterRoutes(r)

	// Registrar controladores e rotas para usuários
	userService := service.NewUsersService(s.db)
	userController := controllers.NewUsersController(userService)
	userController.RegisterRoutes(r)

	return r.Run(configuration.Config.ServerPort)
} // Fim RunServer
