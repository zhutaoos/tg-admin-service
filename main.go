package main

import (
	"app/internal/config"
	"app/internal/model"
	"app/internal/provider"
	"app/internal/router"
	"app/tools/logger"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var err error

func main() {
	// 解析命令行参数
	flag.StringVar(&config.Mode, "mode", "dev", "-mode=prod, -mode=dev")
	flag.StringVar(&config.InitDb, "initDb", "true", "-initDb=true, -initDb=false")
	flag.Parse()

	// 设置时区
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	// 初始化基础配置
	initBasicConfig()

	// 创建Fx应用
	app := fx.New(
		// 提供命令行参数作为依赖
		fx.Provide(
			fx.Annotated{
				Name: "mode",
				Target: func() string {
					return config.Mode
				},
			},
			fx.Annotated{
				Name: "initDb",
				Target: func() string {
					return config.InitDb
				},
			},
		),

		// 所有模块
		provider.AllModules,

		// 启动逻辑
		fx.Invoke(runApplication),
	)

	// 启动应用
	app.Run()
}

// initBasicConfig 初始化基础配置
func initBasicConfig() {
	// 设置服务器名称和PID
	config.ServerName = "tg-admin-service"

	pid := os.Getpid()
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(pid))
	f, _ := os.Create(config.ServerName + ".pid")
	_, err = f.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Printf("进程 PID: %d 写入失败 \n", pid)
		return
	}
	fmt.Printf("进程 PID: %d \n", pid)

	// 获取本地IP
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本地 ip 失败" + err.Error())
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrList {
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				config.LocalIp = ipNet.IP.String()
				break
			}
		}
	}

	// 初始化日志
	logger.Init()
}

// runApplication 运行应用程序
func runApplication(
	lc fx.Lifecycle,
	router *router.Router,
	db *gorm.DB,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 数据库初始化
			if config.InitDb == "true" {
				logger.System("START INIT TABLE ====================")
				// 使用注入的数据库实例进行表初始化
				if err := db.AutoMigrate(&model.User{}, &model.Token{}, &model.Admin{}); err != nil {
					logger.Error("数据库表初始化失败", "error", err)
					return err
				}
				logger.System("END INIT TABLE ====================")
			}

			// 在goroutine中启动服务器，避免阻塞
			go func() {
				if err := router.Run(); err != nil {
					logger.Error("服务器启动失败", "error", err)
				}
			}()

			logger.System("应用程序启动完成")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.System("应用程序正在关闭")
			return nil
		},
	})
}
