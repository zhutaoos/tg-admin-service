# Goingo

```text
 ██████╗   ██████╗  ██╗ ███╗   ██╗  ██████╗   ██████╗ 
██╔════╝  ██╔═══██╗ ██║ ████╗  ██║ ██╔════╝  ██╔═══██╗
██║  ███╗ ██║   ██║ ██║ ██╔██╗ ██║ ██║  ███╗ ██║   ██║
██║   ██║ ██║   ██║ ██║ ██║╚██╗██║ ██║   ██║ ██║   ██║
╚██████╔╝ ╚██████╔╝ ██║ ██║ ╚████║ ╚██████╔╝ ╚██████╔╝
 ╚═════╝   ╚═════╝  ╚═╝ ╚═╝  ╚═══╝  ╚═════╝   ╚═════╝ 
```

基于 Gin + Gorm 整合的开发框架，用于快速构建 API 服务

## 使用技术

- 路由，中间件 [Gin](https://gin-gonic.com/zh-cn/)
- model [Gorm](https://gorm.io/zh_CN/)
- 配置文件解析 [viper](https://github.com/spf13/viper/)
- [jwt](https://github.com/golang-jwt/jwt/)
- [redis](https://redis.uptrace.dev/zh/)

## 目录结构

```
├── config              // 项目配置文件
│   ├── dev.ini         // 开发环境
│   ├── prod.ini        // 测试环境
│   └── deploy.ini      // 打包配置
├── deploy              // 打包上传到正式环境
│   ├── deploy.go
│   ├── deploy.sh
│   └── run.sh
├── internal            // 业务代码
│   ├── logic           // 业务逻辑
│   ├── middleware      // 中间件
│   ├── model           // 模型
│   ├── router          // 路由
│   └── server          // 接口
│   └── global.go       // 存放全局变量和常量
├── log                 // 运行日志
├── main.go             // 入口文件
└── tools               // 通用工具
    ├── conv
    ├── jwt            
    ├── key_utils
    ├── logger          // 日志
    ├── queue           // 队列
    ├── random
    ├── resp            // 响应
    └── utils.go
```

## 运行

```shell
go run main.go -mode=dev
# 运行参数
# -mode=dev    运行测试环境 dev.ini
# -mode=prod   运行正式环境 prod.ini
# -initDb=true 根据结构体初始化数据库
```

## 队列

系统的消息发送链路依托 Redis Stream 与 ZSet：

- 即时任务写入 `tg:send:ready:<chatID>`，延迟任务写入 `tg:send:delayed:<chatID>`；
- `Producer` 会按 chatID 将作业入队，并在首次入队时唤起对应 runner；
- `RunnerManager` 负责维护活跃 chat 列表，按需拉起/回收 mover 与 worker，以控制资源占用；
- `Mover` 将到期的延迟任务搬运到 stream，`Worker` 在单 chat 内顺序消费并调用 Telegram bot，所有限流信息均存储在 Redis 键中。

### Redis 键名约定

- `tg:send:ready:<chatID>`：待发送消息流
- `tg:send:delayed:<chatID>`：延迟任务集合
- `tg:send:cg:<chatID>`：消费组
- `tg:lim:botcnt:<bot>:<sec>`：单 bot 全局速率窗口

### 队列配置项（[queue]）

| 配置项 | 默认值 | 说明 |
| --- | --- | --- |
| `global_rate_per_sec` | 15 | 单 bot 每秒允许的发送次数，0 表示不发送 |
| `mover_batch` | 200 | 每轮搬运的延迟任务数量上限 |
| `mover_interval_ms` | 100 | 搬运器轮询间隔（毫秒） |
| `horizon_sec` | 120 | 用于估算背压阈值的窗口（秒） |
| `stream_max_len` | 0 | Stream 限长，0 表示不限制 |
| `max_chat_runners` | 0 | 允许同时活跃的 chat runner 数量，0 表示不限制 |
| `idle_runner_ttl_ms` | 600000 | chat runner 空闲多久后自动回收（毫秒，<=0 表示不回收） |
| `candidate_cache_ttl_ms` | 3600000 | 候选 bot 列表缓存时间（毫秒） |
| `failure_backoff_ms` | 60000,300000,900000,3600000,10800000,21600000,43200000,86400000 | 单 bot 连续失败退避阶梯（毫秒，逗号分隔） |
| `failure_max_duration_ms` | 86400000 | 单 bot 连续失败允许的最长时间，超过后默认阻断（毫秒） |

根据业务规模调整上述参数，即可在保序与限流之间取得平衡。

## 打包上传到服务器

```shell
cd deploy && go run deploy.go
```
