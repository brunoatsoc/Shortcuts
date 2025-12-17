package controllers

import (
	"intellisys-backend/types"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService types.UsersInterface
}

func NewUsersController(usersService types.UsersInterface) *UsersController {
	return &UsersController{usersService: usersService}
}

func (c *UsersController) RegisterRoutes(r *gin.Engine) {
	r.GET("user", c.getUserByUsername)
}

func (c *UsersController) getUserByUsername(ctx *gin.Context) {
	user := types.UsersRequest{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := c.usersService.GetUserByUsername(user)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, response)
}
