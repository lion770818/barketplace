package main

import (
	"fmt"
	"log"
	"marketplace/internal/common/logs"
	"marketplace/internal/common/redis"
	"marketplace/internal/common/servers"
	"marketplace/internal/common/signals"
	"marketplace/internal/config"
	"marketplace/internal/infrastructure/auth"
	"marketplace/internal/infrastructure/repositories"
	"marketplace/internal/interfaces/web"
	"time"
)

// 建立微服務
func NewServers(cfg *config.SugaredConfig) servers.ServerInterface {

	// 建立redis連線
	redisCfg := &redis.RedisParameter{
		Network:      "tcp",
		Address:      fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  time.Second * time.Duration(10),
		ReadTimeout:  time.Second * time.Duration(10),
		WriteTimeout: time.Second * time.Duration(10),
		PoolSize:     10,
	}
	redisClient, err := redis.NewRedis(redisCfg)
	if err != nil {
		logs.Errorf("newRedis error=%v", err)
	}

	// 建立db連線
	repos, err := repositories.NewRepositories(cfg)
	if err != nil {
		logs.Errorf("newRepositories error=%v", err)
		panic(err)
	}

	// 自動遷移db schema
	repos.Automigrate()

	// 建立redis連線
	redisService, err := auth.NewRedisDB(redisClient.GetClient())
	if err != nil {
		log.Fatal(err)
	}

	apps := web.NewApps(repos, redisService) // db 創建 usercase, 返回usecase
	servers := servers.NewServers()
	servers.AddServer(web.NewWebServer(cfg, apps)) // 啟動 web server

	return servers
}

func main() {

	// 讀取配置文件 .env 或 docker-compose 的 環境變數
	cfg := config.NewConfigFromEnv()
	//log.Printf("cfg=%v", cfg)

	// 日誌設定
	logConfig := logs.LogConfig{
		Env:        cfg.Log.Env,
		Path:       cfg.Log.Path,
		Encoding:   cfg.Log.Encoding,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
	}
	// 初始化日志
	logs.Init(logConfig)
	logs.Debugf("service start......................")

	// 建立service微服務
	servers := NewServers(cfg)

	// 啟動 servers
	servers.AsyncStart()

	logs.Debugf("優雅關閉 等待訊號中...")
	signals.WaitWith(servers.Stop)

	logs.Debugf("service end......................")

}
