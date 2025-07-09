package router

import (
	"github.com/gin-gonic/gin"
)

type IndexRoute struct {
	group *gin.RouterGroup
}

func (r *IndexRoute) initRoute() {
	r.group.POST("login", index_api.Login) // 登陆
}
