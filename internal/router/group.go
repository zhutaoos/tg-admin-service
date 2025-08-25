package router

import (
	group_controller "app/internal/controller/group"

	"github.com/gin-gonic/gin"
)

type GroupRoute struct {
	groupController *group_controller.GroupController
}

// NewGroupRoute 创建群组路由
func NewGroupRoute(groupController *group_controller.GroupController) *GroupRoute {
	return &GroupRoute{
		groupController: groupController,
	}
}

// InitRoute 初始化群组路由
func (gr *GroupRoute) InitRoute(router *gin.Engine) {
	// 群组管理路由
	group := router.Group("/api/group")
	{
		group.POST("/create", gr.groupController.CreateGroup)
		group.POST("/update", gr.groupController.UpdateGroup)
		group.POST("/delete", gr.groupController.DeleteGroup)
		group.POST("/list", gr.groupController.SearchGroups)
		group.GET("/my", gr.groupController.GetMyGroups)
	}
}
