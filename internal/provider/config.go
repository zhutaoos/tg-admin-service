package provider

import (
	"app/internal/config"
	"app/tools/logger"
	"context"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/fx"
)

type ConfigParams struct {
	fx.In
	Mode   string `name:"mode"`
	InitDb string `name:"initDb"`
}

// NewConfig 创建配置实例
func NewConfig(params ConfigParams) (*config.Config, error) {
	conf := &config.Config{
		Path:     "./config",
		FileName: params.Mode, // dev or prod
	}
	return conf.Init(), nil
}

// NewDatabaseConfig 创建数据库配置
func NewDatabaseConfig(conf *config.Config) *config.DbConf {
	return &config.DbConf{
		UserName: config.Get[string](conf, "mysql", "username"),
		Password: config.Get[string](conf, "mysql", "password"),
		Ip:       config.Get[string](conf, "mysql", "ip"),
		Port:     config.Get[string](conf, "mysql", "port"),
		DbName:   config.Get[string](conf, "mysql", "db_name"),
	}
}

// NewRedisConfig 创建Redis配置
func NewRedisConfig(conf *config.Config) *config.RedisConf {
	return &config.RedisConf{
		Ip:       config.Get[string](conf, "redis", "ip"),
		Port:     config.Get[string](conf, "redis", "port"),
		Username: config.Get[string](conf, "redis", "username"),
		Password: config.Get[string](conf, "redis", "password"),
		Db:       config.Get[int](conf, "redis", "db"),
		MaxTotal: config.Get[int](conf, "redis", "max_total"),
	}
}

// ConfigWatcher 配置文件监听器
type ConfigWatcher struct {
	config  *config.Config
	watcher *fsnotify.Watcher
}

// NewConfigWatcher 创建配置监听器，实现热重载
func NewConfigWatcher(lc fx.Lifecycle, conf *config.Config, params ConfigParams) *ConfigWatcher {
	watcher := &ConfigWatcher{config: conf}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return watcher.Start(params.Mode)
		},
		OnStop: func(ctx context.Context) error {
			return watcher.Stop()
		},
	})

	return watcher
}

// Start 启动配置文件监听
func (cw *ConfigWatcher) Start(mode string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	cw.watcher = watcher

	configFile := filepath.Join("./config", mode+".ini")
	err = watcher.Add(configFile)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logger.System("配置文件发生变化，正在重新加载", "file", event.Name)
					// 这里可以添加配置重新加载逻辑
					// 由于config.Config的限制，暂时只记录日志
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("配置文件监听错误", "error", err)
			}
		}
	}()

	logger.System("配置文件监听器已启动", "file", configFile)
	return nil
}

// Stop 停止配置文件监听
func (cw *ConfigWatcher) Stop() error {
	if cw.watcher != nil {
		logger.System("正在停止配置文件监听器")
		return cw.watcher.Close()
	}
	return nil
}
