package api

import (
	"database/sql"
	"intellisys-backend/configuration"
	"intellisys-backend/controllers"
	"intellisys-backend/service"

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

	shortcutService := service.NewShortcutService(s.db)
	shortcutController := controllers.NewShortcutController(shortcutService)
	shortcutController.RegisterRoutes(r)

	return r.Run(configuration.Config.ServerPort)
}
