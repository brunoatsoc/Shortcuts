package controllers

import (
	"intellisys-backend/internal/middlewares"
	"intellisys-backend/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShortcutController struct {
	service types.ShortcutsInterface
}

func NewShortcutController(service types.ShortcutsInterface) *ShortcutController {
	return &ShortcutController{service: service}
}

// Método para registrar as rotas de atalho
func (c *ShortcutController) RegisterRoutes(g *gin.Engine) {
	g.GET("shortcut", c.getAllShortcuts)

	auth := g.Group("/")

	// Aplicar middleware de autenticação nas rotas que são protegidas
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("shortcut", c.createShortcut)
		auth.PUT("shortcut/:id", c.updateShortcut)
		auth.DELETE("shortcut/:id", c.deleteShortcut)
	}
} // Fim RegisterRoutes

// Método para criar um novo atalho
func (c *ShortcutController) createShortcut(ctx *gin.Context) {
	var req types.ShortcutsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	res, err := c.service.CreateShortcut(req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create shortcut"})
		return
	}

	ctx.JSON(201, res)
} // Fim CreateShortcut

// Método para obter todos os atalhos
func (c *ShortcutController) getAllShortcuts(ctx *gin.Context) {
	res, err := c.service.GetAllShortcuts()

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve shortcuts"})
		return
	}

	ctx.JSON(200, res)
} // Fim GetAllShortcuts

// Método para atualizar um atalho existente
func (c *ShortcutController) updateShortcut(ctx *gin.Context) {
	var req types.ShortcutsRequest
	idParam := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	idInt, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID parameter"})
		return
	}

	res, err := c.service.UpdateShortcut(idInt, req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update shortcut"})
		return
	}

	ctx.JSON(200, res)
} // Fim UpdateShortcut

// Método para deletar um atalho
func (c *ShortcutController) deleteShortcut(ctx *gin.Context) {
	idParam := ctx.Param("id")

	idInt, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := c.service.DeleteShortcut(idInt); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete shortcut"})
		return
	}

	ctx.JSON(204, nil)
} // Fim DeleteShortcut
