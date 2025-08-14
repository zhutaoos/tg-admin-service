package user

import (
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) UserList(ctx *gin.Context) {
	var req request.UserSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失"}).Response()
	}

	// 调用业务层获取用户列表
	users, total, err := uc.ListUser(req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询用户列表失败: " + err.Error()}).Response()
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

func (uc *UserController) GetUserInfo(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户ID不能为空"}).Response()
	}

	user, err := uc.LoadUser(userId)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "用户不存在"}).Response()
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取用户信息成功", Data: user}).Response()
}
