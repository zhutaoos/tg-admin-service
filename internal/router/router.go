package router

import (
	"app/internal/config"
	"app/internal/middleware"
	"app/internal/service"
	"app/tools/logger"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Router 主路由结构
type Router struct {
	Engine        *gin.Engine
	AdminRoute    *AdminRoute
	UserRoute     *UserRoute
	IndexRoute    *IndexRoute
	EvaluateRoute *EvaluateRoute
	Config        *config.Config
	TokenService  service.TokenService
	adminService  service.AdminService
}

// NewRouter 创建路由实例
func NewRouter(
	adminRoute *AdminRoute,
	userRoute *UserRoute,
	indexRoute *IndexRoute,
	evaluateRoute *EvaluateRoute,
	conf *config.Config,
	tokenService service.TokenService,
	adminService service.AdminService,
) *Router {
	return &Router{
		AdminRoute:    adminRoute,
		UserRoute:     userRoute,
		IndexRoute:    indexRoute,
		EvaluateRoute: evaluateRoute,
		Config:        conf,
		TokenService:  tokenService,
		adminService:  adminService,
	}
}

// SetupEngine 设置Gin引擎
func (router *Router) SetupEngine() *gin.Engine {
	// 设置Gin模式
	if config.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 设置受信任的代理
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	// 设置日志
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

	// 响应中间件
	r.Use(middleware.RespMiddleware())

	// JWT中间件白名单
	whitelist := []string{
		"/api/admin/login",   // 管理员登录
		"/api/admin/initPwd", // 初始化密码
		"/api/index/health",  // 健康检查
	}
	r.Use(middleware.JwtMiddlewareWithWhitelist(whitelist, router.TokenService, router.adminService))

	router.Engine = r
	return r
}

// InitRoutes 初始化所有路由
func (router *Router) InitRoutes() {
	// 初始化各个模块的路由
	router.AdminRoute.InitRoute(router.Engine)
	router.UserRoute.InitRoute(router.Engine)
	router.IndexRoute.InitRoute(router.Engine)
	router.EvaluateRoute.InitRoute(router.Engine)
}

// Run 启动服务器
func (router *Router) Run() error {
	// 设置引擎
	router.SetupEngine()

	// 初始化路由
	router.InitRoutes()

	// 获取端口配置
	port := config.Get[string](router.Config, "server", "port")

	// 启动服务器
	logger.System("服务器启动", "port", port)
	return router.Engine.Run(":" + port)
}
