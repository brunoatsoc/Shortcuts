package controller

import (
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/user"
	"intellisys-backend/internal/infra/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *user.AuthService
	expireTime  int64
}

func NewAuthController(authService *user.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		expireTime:  cfg.JWT.ExpireTime,
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
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"access_token",
		token,
		int(a.expireTime)*60,
		"/",
		"",
		false, // Secure: só HTTPS
		true,  // HttpOnly: JS não acessa
	)
}
