package admin_controller

import (
	"app/internal/config"
	"app/internal/model"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminController 管理员控制器
type AdminController struct {
	adminService service.AdminService
}

// NewAdminController 创建管理员控制器实例
func NewAdminController(adminService service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// AdminLogin 管理员登录
func (ac *AdminController) AdminLogin(ctx *gin.Context) {
	var req request.AdminLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失"}).Response()
	}

	// 调用业务层处理登录
	result, err := ac.adminService.Login(req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: err.Error()}).Response()
	}

	// 组装响应数据
	data := map[string]interface{}{
		"token":      result.Token,
		"token_info": result.TokenInfo,
		"user":       result.User,
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "登录成功", Data: data}).Response()
}

// InitPwd 初始化密码
func (ac *AdminController) InitPwd(ctx *gin.Context) {
	var req request.InitPwdRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "请输入密码"}).Response()
	}

	if req.Password == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "请输入密码"}).Response()
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "密码加密失败"}).Response()
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "密码初始化成功", Data: string(hashedPassword)}).Response()
}

// Profile 获取用户信息
func (ac *AdminController) Profile(ctx *gin.Context) {
	user := ctx.MustGet(config.CurrentUser).(*model.Admin)

	resp.Ok(user)
}
