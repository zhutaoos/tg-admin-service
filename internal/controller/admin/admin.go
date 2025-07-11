package admin_api

import (
	"app/internal/logic"
	"app/internal/model"
	"app/internal/request"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin 管理员登录
func AdminLogin(content *gin.Context) {
	var req request.AdminLoginRequest
	if err := content.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失"}).Response()
	}

	admin := &model.Admin{
		Account: req.Username,
	}
	admin = admin.GetAdmin()

	// 2. 检查用户是否存在
	if admin.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "账号不存在", Data: nil}).Response()
	}

	// 3. 使用 bcrypt 比较密码（重要！）
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "密码错误", Data: nil}).Response()
	}

	// 4. 密码正确，生成JWT令牌
	data := make(map[string]interface{})

	j, userJwt := logic.TokenLogicInstance.GenerateJwt(admin.Id, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = admin
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "登陆成功", Data: data}).Response()
}

func InitPwd(content *gin.Context) {
	var req request.InitPwdRequest
	if err := content.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "请输入密码"}).Response()
	}

	if req.Password == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "请输入密码"}).Response()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "密码错误", Data: nil}).Response()
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "密码初始化成功", Data: string(hashedPassword)}).Response()
}

func Profile(content *gin.Context) {
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取用户信息成功", Data: nil}).Response()
}
