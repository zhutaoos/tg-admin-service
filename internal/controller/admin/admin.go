package admin_api

import (
	"app/internal/logic"
	model2 "app/internal/model"
	"app/tools"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminLogin 管理员登录
func AdminLogin(content *gin.Context) {
	account := content.PostForm("account")
	password := content.PostForm("password")

	if account == "" || password == "" {
		(&resp.JsonResp{Code: resp.ReFail, Message: "请输入账号密码", Body: nil}).Response()
	}

	admin := &model2.Admin{
		Account:  account,
		Password: tools.Md5(password, model2.UserPwdSalt),
	}

	admin = admin.GetAdmin()
	if admin.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "账号密码错误", Body: nil}).Response()
	}

	data := make(map[string]interface{})
	j, userJwt := logic.TokenLogicInstance.GenerateJwt(admin.Id, jwt.AdminJwtType, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = admin
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: data}).Response()
}

// GetAdminInfo 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	adminId := c.GetUint(string(jwt.AdminJwtType))
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: logic.GetAdminAuth(adminId)}).Response()
}

// GetAdminList 获取管理员列表
func GetAdminList(_ *gin.Context) {
	admin := model2.Admin{}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: admin.GetList(0)}).Response()
}

// DelAdmin 删除管理员
func DelAdmin(c *gin.Context) {
	ids := c.PostForm("ids")
	l := strings.Split(ids, ",")
	var temp []int
	for _, v := range l {
		i, err := conv.Conv[int](v)
		if err == nil {
			temp = append(temp, i)
		}
	}
	article := &model2.Admin{}
	article.DelAdmin(temp)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: nil}).Response()
}
