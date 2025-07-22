package router

import (
	"app/internal/config"
	"app/internal/middleware"
	"app/tools/logger"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter(port string) {
	var err error

	if config.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.CORSMiddleware()) // 解决跨域

	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	f, _ := os.OpenFile(logger.AccessLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	c := gin.LoggerConfig{
		Output:    f,
		SkipPaths: []string{"/favicon.ico"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\" POSTFORM - [%s] \n",
				params.ClientIP,
				params.TimeStamp.Format(time.DateTime),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
				params.Request.PostForm,
			)
		},
	}
	r.Use(gin.LoggerWithConfig(c))
	r.Use(middleware.RespMiddleware()) // 响应中间件

	whitelist := []string{
		"/api/admin/login",   // 管理员登录
		"/api/admin/initPwd", // 初始化密码
		"/api/health",        // 健康检查
	}
	r.Use(middleware.JwtMiddlewareWithWhitelist(whitelist))

	AdminRoute := AdminRoute{group: r.Group("api/admin")}
	IndexRoute := IndexRoute{group: r.Group("api/index")}
	UserRoute := UserRoute{group: r.Group("api/user")}
	AdminRoute.initRoute()
	IndexRoute.initRoute()
	UserRoute.initRoute()

	err = r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
}
