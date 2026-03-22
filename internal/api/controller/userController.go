package controller

import (
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *user.UserService
}

func NewUserController(service *user.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) UpdateUsuario(ctx *gin.Context) {
	loggedUserId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	var req dto.UserRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateUser(ctx, req, loggedUserId.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) GetUserById(ctx *gin.Context) {
	loggedUserId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	result, err := c.service.GetUserById(ctx, loggedUserId.(string))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *UserController) GetAll(ctx *gin.Context) {
	_, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	results, err := c.service.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	loggedUserId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	if err := c.service.DeleteUser(ctx, loggedUserId.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
