package controllers

import (
	"intellisys-backend/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShortcutController struct {
	service types.ShortcutsInterface
}

func NewShortcutController(service types.ShortcutsInterface) *ShortcutController {
	return &ShortcutController{service: service}
}

func (c *ShortcutController) RegisterRoutes(g *gin.Engine) {
	g.POST("shortcut", c.createShortcut)
	g.GET("shortcut", c.getAllShortcuts)
	g.PUT("shortcut/:id", c.updateShortcut)
	g.DELETE("shortcut/:id", c.deleteShortcut)
}

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
}

func (c *ShortcutController) getAllShortcuts(ctx *gin.Context) {
	res, err := c.service.GetAllShortcuts()

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve shortcuts"})
		return
	}

	ctx.JSON(200, res)
}

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
}

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
}
