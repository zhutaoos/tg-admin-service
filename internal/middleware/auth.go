package middleware

import (
	logic2 "app/internal/logic"
	"app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckJwt() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "请上传jwt", Data: nil}).Response()
		}
		data, err := logic2.TokenLogicInstance.CheckJwt(token)
		if err != nil {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "jwt解析失败", Data: map[string]any{"err": err.Error()}}).Response()
		}

		switch data.Type {
		case jwt.AdminJwtType:
			user := &model.Admin{
				Id: data.Uid,
			}
			user = user.GetAdmin()
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "未查询到用户", Data: nil}).Response()
			}
			isSuper := user.IsSuper == 1

			rolesGroup := new(model.RolesGroup)
			rolesGroup.Id = user.RolesGroupId
			rolesGroup.GetRolesGroup()
			auth := logic2.NewAdminAuth(user.Id, user.Pid, rolesGroup, isSuper)
			auth.Name = user.Name
			auth.Avatar = user.Avatar
			auth.Cache()
			c.Set(string(jwt.AdminJwtType), auth.Id) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		case jwt.IndexJwtType:
			user := logic2.UserLogicInstance.LoadUser(data.Uid)
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "未查询到用户", Data: nil}).Response()
			}
			c.Set(string(jwt.IndexJwtType), user)
			c.Next()
			break
		}

	}
}

// BackendAuth 管理后台鉴权
func BackendAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminId := c.GetUint(string(jwt.AdminJwtType))
		admin := logic2.GetAdminAuth(adminId)
		fmt.Println(admin)
		fmt.Println(c.Request.URL.Path)
		menu := model.Roles{}
		menu.SearchByPath(c.Request.URL.Path)
		if menu.Id == 0 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "BackendAuth 未查询到权限", Data: nil}).Response()
		}
		checkId, _ := conv.Conv[uint](menu.Id)
		has, _ := conv.InSlice[uint](admin.RolesIds, checkId)
		if has == -1 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "BackendAuth 无权限访问", Data: nil}).Response()
		}
		c.Next()
	}
}

// JwtMiddlewareWithWhitelist JWT中间件（支持白名单）
func JwtMiddlewareWithWhitelist(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查当前请求路径是否在白名单中
		currentPath := c.Request.URL.Path
		for _, path := range whitelist {
			// 支持精确匹配和前缀匹配
			if currentPath == path || strings.HasPrefix(currentPath, path) {
				c.Next()
				return
			}
		}

		// 不在白名单中，执行JWT鉴权逻辑
		token := c.Request.Header.Get("Token")
		if token == "" {
			resp.NeedLogin().Response()
		}

		data, err := logic2.TokenLogicInstance.CheckJwt(token)
		if err != nil {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "jwt解析失败"}).Response()
		}

		switch data.Type {
		case jwt.AdminJwtType:
			user := &model.Admin{
				Id: data.Uid,
			}
			user = user.GetAdmin()
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "未查询到用户"}).Response()
			}
			isSuper := user.IsSuper == 1

			rolesGroup := new(model.RolesGroup)
			rolesGroup.Id = user.RolesGroupId
			rolesGroup.GetRolesGroup()
			auth := logic2.NewAdminAuth(user.Id, user.Pid, rolesGroup, isSuper)
			auth.Name = user.Name
			auth.Avatar = user.Avatar
			auth.Cache()
			c.Set(string(jwt.AdminJwtType), auth.Id) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		case jwt.IndexJwtType:
			user := logic2.UserLogicInstance.LoadUser(data.Uid)
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "未查询到用户"}).Response()
			}
			c.Set(string(jwt.IndexJwtType), user)
			c.Next()
			break
		}
	}
}
