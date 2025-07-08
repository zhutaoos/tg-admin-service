package main

import (
	"app/internal/global"
	"app/internal/model"
	"app/internal/router"
	config "app/tools/config"
	"app/tools/logger"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var err error

func main() {
	flag.StringVar(&global.Mode, "mode", "dev", "-mode=prod, -mode=dev") // "dev" or "prod"
	flag.StringVar(&global.InitDb, "initDb", "false", "-initDb=true, -initDb=false")
	flag.Parse()
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	conf := (&config.Config{
		Path:     "./config",
		FileName: global.Mode, // dev or prod
	}).Init()

	global.ServerName = config.Get[string](conf, "server", "name")
	global.Version = config.Get[string](conf, "server", "version")
	pid := os.Getpid()
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(pid))
	f, _ := os.Create(global.ServerName + ".pid")
	_, err = f.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Printf("进程 PID: %d 写入失败 \n", pid)
		return
	}
	fmt.Printf("进程 PID: %d \n", pid)

	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本地 ip 失败" + err.Error())
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrList {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				global.LocalIp = ipNet.IP.String()
				break
			}
		}
	}

	logger.Init()

	model.InitDb(&model.DbConf{
		UserName: config.Get[string](conf, "mysql", "username"),
		Password: config.Get[string](conf, "mysql", "password"),
		Ip:       config.Get[string](conf, "mysql", "ip"),
		Port:     config.Get[string](conf, "mysql", "port"),
		DbName:   config.Get[string](conf, "mysql", "db_name"),
	})

	model.InitRedis(&model.RedisConf{
		Ip:   config.Get[string](conf, "redis", "ip"),
		Port: config.Get[string](conf, "redis", "port"),
	})

	if global.InitDb == "true" {
		logger.System("START INIT TABLE ====================")
		m := new(model.MysqlBaseModel)
		m.SetTableComment("用户表").CreateTable(model.User{})
		m.SetTableComment("token").CreateTable(model.Token{})
		m.SetTableComment("article").CreateTable(model.Article{})
		m.SetTableComment("").CreateTable(model.Admin{})
		m.SetTableComment("").CreateTable(model.Cate{})
		m.SetTableComment("").CreateTable(model.Roles{})
		m.SetTableComment("角色表").CreateTable(model.RolesGroup{})
		logger.System("END INIT TABLE ====================")
	}

	port := config.Get[string](conf, "server", "port")
	router.InitRouter(port)
}
