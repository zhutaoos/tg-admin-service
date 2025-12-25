package controller

import (
	"app/internal/config"
	"app/internal/model"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) *model.Admin {
	if user, exists := c.Get(config.CurrentUser); exists {
		if admin, ok := user.(*model.Admin); ok {
			return admin
		}
	}
	panic("用户未登录")
}

func CurrentUserId(c *gin.Context) uint {
	userId := c.GetUint(config.CurrentUserId)
	if userId > 0 {
		return userId
	}
	panic("用户未登录")
}
