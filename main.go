package main

import (
	"app/internal/config"
	"app/internal/model"
	"app/internal/router"
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
	flag.StringVar(&config.Mode, "mode", "dev", "-mode=prod, -mode=dev") // "dev" or "prod"
	flag.StringVar(&config.InitDb, "initDb", "true", "-initDb=true, -initDb=false")
	flag.Parse()
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	conf := (&config.Config{
		Path:     "./config",
		FileName: config.Mode, // dev or prod
	}).Init()

	config.ServerName = config.Get[string](conf, "server", "name")
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
				config.LocalIp = ipNet.IP.String()
				break
			}
		}
	}

	logger.Init()

	config.InitMysql(&config.DbConf{
		UserName: config.Get[string](conf, "mysql", "username"),
		Password: config.Get[string](conf, "mysql", "password"),
		Ip:       config.Get[string](conf, "mysql", "ip"),
		Port:     config.Get[string](conf, "mysql", "port"),
		DbName:   config.Get[string](conf, "mysql", "db_name"),
	})

	// config.InitRedis(&config.RedisConf{
	// 	Ip:       config.Get[string](conf, "redis", "ip"),
	// 	Port:     config.Get[string](conf, "redis", "port"),
	// 	Username: config.Get[string](conf, "redis", "username"),
	// 	Password: config.Get[string](conf, "redis", "password"),
	// 	Db:       config.Get[int](conf, "redis", "db"),
	// 	MaxTotal: config.Get[int](conf, "redis", "max_total"),
	// })

	if config.InitDb == "true" {
		logger.System("START INIT TABLE ====================")
		m := new(model.MysqlBaseModel)
		m.SetTableComment("用户表").CreateTable(model.User{})
		m.SetTableComment("token").CreateTable(model.Token{})
		m.SetTableComment("").CreateTable(model.Admin{})
		logger.System("END INIT TABLE ====================")
	}

	port := config.Get[string](conf, "server", "port")
	router.InitRouter(port)
}
