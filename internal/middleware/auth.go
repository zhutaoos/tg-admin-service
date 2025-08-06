package middleware

import (
	"app/internal/config"
	"app/internal/model"
	"app/internal/service"
	"app/tools/resp"
	"strings"

	"github.com/gin-gonic/gin"
)

// JwtMiddlewareWithWhitelist JWT中间件（支持白名单）
func JwtMiddlewareWithWhitelist(whitelist []string, tokenService service.TokenService, adminService service.AdminService) gin.HandlerFunc {
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

		data, err := tokenService.CheckJwt(token)
		if err != nil {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "jwt解析失败"}).Response()
		}

		// 验证用户是否存在
		user := &model.Admin{}
		err = adminService.GetAdminById(int64(data.UserId), user)
		if err != nil || user.Id <= 0 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Msg: "未查询到用户"}).Response()
		}

		// 将用户信息存储到上下文中，供后续处理使用
		c.Set(config.CurrentUser, user)
		c.Set(config.CurrentUserId, user.Id)

		c.Next()
	}
}
