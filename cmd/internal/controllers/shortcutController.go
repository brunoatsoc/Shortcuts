package controllers

import (
	"intellisys-backend/internal/middlewares"
	"intellisys-backend/internal/types"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortcutController struct {
	service types.ShortcutsInterface
}

func NewShortcutController(service types.ShortcutsInterface) *ShortcutController {
	return &ShortcutController{service: service}
}

// Método para registrar as rotas de atalho
func (c *ShortcutController) RegisterRoutes(g *gin.Engine) {
	// Torna a pasta "./uploads" acessível na url "/uploads"
	g.Static("/uploads", "./uploads")

	g.GET("shortcut", c.getAllShortcuts)

	auth := g.Group("/")

	// Aplica middleware de autenticação nas rotas que são protegidas
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("shortcut", c.createShortcut)
		auth.PUT("shortcut/:id", c.updateShortcut)
		auth.DELETE("shortcut/:id", c.deleteShortcut)
	}
} // Fim RegisterRoutes

// Metorodo para criar um novo atalho
func (c *ShortcutController) createShortcut(ctx *gin.Context) {
	var req types.ShortcutsRequest

	// Tenta fazer o bind dos campos de texto (Name, URL)
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Processa o upload da imagem
	imagePath, err := c.saveFile(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upload image"})
		return
	}

	if imagePath != "" {
		req.ImageURL = imagePath
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

// Metodo para atualizar um atalho
func (c *ShortcutController) updateShortcut(ctx *gin.Context) {
	var req types.ShortcutsRequest
	idParam := ctx.Param("id")
	idInt, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Verifica se enviou uma nova imagem
	newImagePath, err := c.saveFile(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upload image"})
		return
	}

	// Se enviou nova imagem, atualiza.
	if newImagePath != "" {
		req.ImageURL = newImagePath
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

// Helper para salvar arquivo
func (c *ShortcutController) saveFile(ctx *gin.Context) (string, error) {
	// Tenta pegar o arquivo com a chave "image" do form-data
	file, err := ctx.FormFile("image")

	// Se não houver arquivo (err != nil), retorna string vazia sem erro
	if err != nil {
		return "", nil
	}

	// Gera um nome único para não sobrescrever arquivos
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	path := "uploads/" + filename

	// Salva o arquivo na pasta uploads
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	// Retorna o caminho com uma barra na frente para ser acessível via URL
	return "/" + path, nil
} // Fim saveFile
