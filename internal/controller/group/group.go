package group

import (
	"app/internal/controller"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

// GroupController 群组管理控制器
type GroupController struct {
	controller.BaseController
	groupService service.GroupService
}

// NewGroupController 创建群组控制器实例
func NewGroupController(groupService service.GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

func (c *GroupController) CreateGroup(ctx *gin.Context) {
	var req request.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	currentUserId := c.CurrentUserId(ctx)
	err := c.groupService.CreateGroup(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建群组关联失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建成功"}).Response()
}

func (c *GroupController) UpdateGroup(ctx *gin.Context) {
	var req request.UpdateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	currentUserId := c.CurrentUserId(ctx)
	if err := c.groupService.UpdateGroup(ctx, req, currentUserId); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "更新群组信息失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新成功"}).Response()
}

func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	var req request.DeleteGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	currentUserId := c.CurrentUserId(ctx)

	if err := c.groupService.DeleteGroup(ctx, req, currentUserId); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "删除群组关联失败: " + err.Error()}).Response()
		return
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "删除成功"}).Response()
}

func (c *GroupController) SearchGroups(ctx *gin.Context) {
	var req request.SearchGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	currentUserId := c.CurrentUserId(ctx)
	groups, total, err := c.groupService.SearchGroups(ctx, req, currentUserId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "查询群组列表失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "查询成功",
		Data: map[string]interface{}{
			"total": total,
			"list":  groups,
		},
	}).Response()
}

func (c *GroupController) GetMyGroups(ctx *gin.Context) {
	currentUserId := c.CurrentUserId(ctx)

	groups, err := c.groupService.GetGroupsByAdminID(ctx, int(currentUserId))
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "获取群组列表失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "获取成功",
		Data: groups,
	}).Response()
}