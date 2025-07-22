package router

import (
	userApi "app/internal/controller/user"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	group *gin.RouterGroup
}

func (r *UserRoute) initRoute() {
	// 当前实现：每个路由单独配置JWT中间件
	r.group.POST("userList", userApi.UserList)

}
