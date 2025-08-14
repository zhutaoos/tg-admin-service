package router

import (
	admin_controller "app/internal/controller/admin"

	"github.com/gin-gonic/gin"
)

type GroupRoute struct {
	groupController *admin_controller.GroupController
}

// NewGroupRoute 创建群组路由
func NewGroupRoute(groupController *admin_controller.GroupController) *GroupRoute {
	return &GroupRoute{
		groupController: groupController,
	}
}

// InitRoute 初始化群组路由
func (gr *GroupRoute) InitRoute(adminGroup *gin.RouterGroup) {
	// 群组管理路由
	group := adminGroup.Group("/group")
	{
		group.POST("/create", gr.groupController.CreateGroup)
		group.POST("/update", gr.groupController.UpdateGroup)
		group.POST("/delete", gr.groupController.DeleteGroup)
		group.POST("/list", gr.groupController.SearchGroups)
		group.GET("/my", gr.groupController.GetMyGroups)
	}
}