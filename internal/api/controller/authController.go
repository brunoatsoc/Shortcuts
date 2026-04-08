package controller

import (
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *user.AuthService
}

func NewAuthController(authService *user.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nome é obrigatório"})
		return
	}

	if req.Password == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "senha é obrigatória"})
		return
	}

	loggedUser, token, err := a.authService.Authenticate(ctx, *req.Name, *req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	a.setAccessTokenCookie(ctx, *token)

	ctx.JSON(http.StatusOK, loggedUser)
}

func (a *AuthController) setAccessTokenCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteLaxMode) // 👈 deve vir antes do SetCookie
	ctx.SetCookie(
		"access_token",
		token,
		15*60, // 15 minutos em segundos
		"/",
		"",
		false, // Secure: só HTTPS
		true,  // HttpOnly: JS não acessa
	)
}
