package admin_api

import (
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"
	"time"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// UserList 用户列表
func (uc *UserController) UserList(ctx *gin.Context) {
	var req request.UserSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失"}).Response()
	}

	// 调用业务层获取用户列表
	users, total, err := uc.userService.UserList(req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询用户列表失败: " + err.Error()}).Response()
	}

	// 处理时间格式化
	for i := range users {
		if users[i].CreateTime > 0 {
			users[i].CreateTimeStr = time.Unix(users[i].CreateTime, 0).Format("2006-01-02 15:04:05")
		}
	}

	// 构建响应数据
	data := map[string]interface{}{
		"list":  users,
		"total": total,
		"page":  req.Page,
		"limit": req.Limit,
		"pages": (total + int64(req.Limit) - 1) / int64(req.Limit), // 总页数
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取用户列表成功", Data: data}).Response()
}

// GetUserInfo 获取用户信息
func (uc *UserController) GetUserInfo(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户ID不能为空"}).Response()
	}

	user, err := uc.userService.LoadUser(userId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户不存在"}).Response()
	}

	// 格式化时间
	if user.CreateTime > 0 {
		user.CreateTimeStr = time.Unix(user.CreateTime, 0).Format("2006-01-02 15:04:05")
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取用户信息成功", Data: user}).Response()
}
