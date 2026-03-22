package controller

import (
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/shortcut"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortcutController struct {
	service *shortcut.ShortcutService
}

func NewshortcutController(service *shortcut.ShortcutService) *ShortcutController {
	return &ShortcutController{
		service: service,
	}
}

func (c *ShortcutController) Create(ctx *gin.Context) {
	var req dto.ShortcutRequestDTO

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := req.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "erro ao abrir imagem",
		})
		return
	}
	defer file.Close()

	input := dto.ShortcutCreateDTO{
		Name:         req.Name,
		ShortcutType: req.ShortcutType,
		Image:        file,
	}

	result, err := c.service.CreateShortcut(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c *ShortcutController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.ShortcutRequestUpdateDTO

	// Faz bind do multipart (mesmo sem imagem)
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var imageFile io.ReadSeeker

	// Se o usuário enviou imagem
	if req.Image != nil {
		file, err := req.Image.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "erro ao abrir imagem",
			})
			return
		}
		defer file.Close()

		imageFile = file
	}

	// Converte para DTO do service
	input := dto.ShortcutUpdateDTO{
		Name:  req.Name,
		Image: imageFile, // pode ser nil
	}

	// Chama service
	err := c.service.UpdateShortcut(ctx, input, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ShortcutController) GetAll(ctx *gin.Context) {
	result, err := c.service.GetAllShortcuts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ShortcutController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteShortcut(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
