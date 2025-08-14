package admin

import (
	"app/internal/controller"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminController 管理员控制器
type AdminController struct {
	controller.BaseController
	adminService service.AdminService
}

// NewAdminController 创建管理员控制器实例
func NewAdminController(adminService service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

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

func (ac *AdminController) Profile(ctx *gin.Context) {
	user := ac.CurrentUser(ctx)
	resp.Ok(user)
}
