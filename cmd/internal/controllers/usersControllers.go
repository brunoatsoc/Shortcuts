package controllers

import (
	"fmt"
	"intellisys-backend/internal/configuration"
	"intellisys-backend/internal/service"
	"intellisys-backend/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService types.UsersInterface
}

func NewUsersController(usersService types.UsersInterface) *UsersController {
	return &UsersController{usersService: usersService}
}

func (c *UsersController) RegisterRoutes(r *gin.Engine) {
	r.POST("/login", c.login)
}

func (c *UsersController) login(ctx *gin.Context) {
	var user types.UsersRequest

	fmt.Println(configuration.Config)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	fmt.Println(user)

	dbUser, err := c.usersService.GetUserByUsername(user)

	fmt.Println(dbUser)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid User or Password"})
		return
	}

	if !service.CheckPassword(user.Password, dbUser.Password) {
		ctx.JSON(401, gin.H{"error": "Invalid User or Password"})
		return
	}

	token, err := service.GenerateToken(strconv.Itoa(dbUser.ID), dbUser.Username)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(200, gin.H{
		"token": token,
	})
}
