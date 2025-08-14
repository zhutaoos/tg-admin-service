package controller

import (
	"app/internal/config"
	"app/internal/model"

	"github.com/gin-gonic/gin"
)

// BaseController 基础控制器，提供通用功能
// 所有控制器都应该嵌入此结构体来获取基础功能
type BaseController struct{}

func (bc *BaseController) CurrentUser(c *gin.Context) *model.Admin {
	if user, exists := c.Get(config.CurrentUser); exists {
		if admin, ok := user.(*model.Admin); ok {
			return admin
		}
	}
	panic("用户未登录")
}

func (bc *BaseController) CurrentUserId(c *gin.Context) uint {
	userId := c.GetUint(config.CurrentUserId)
	if userId > 0 {
		return userId
	}
	panic("用户未登录")
}
