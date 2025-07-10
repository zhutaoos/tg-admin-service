package admin_api

import (
	"app/internal/logic"
	"app/internal/model"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin 管理员登录
func AdminLogin(content *gin.Context) {
	account := content.PostForm("account")
	password := content.PostForm("password")

	if account == "" || password == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "请输入账号密码"}).Response()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "密码错误", Data: nil}).Response()
	}

	admin := &model.Admin{
		Account:  account,
		Password: string(hashedPassword),
	}

	admin = admin.GetAdmin()
	if admin.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "账号密码错误", Data: nil}).Response()
	}

	data := make(map[string]interface{})

	j, userJwt := logic.TokenLogicInstance.GenerateJwt(admin.Id, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = admin
	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "登陆成功", Data: data}).Response()
}
