package config

import (
	"fmt"
	"log"
	"marketplace/internal/common/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SugaredConfig 将配置文件的参数解析,比如解析时间为 time.Ticker
type SugaredConfig struct {
	*Config
	AuthExpireTime time.Duration
}

func NewConfig(filePath string) *SugaredConfig {
	// 初始化配置文件
	pflag.StringP("config", "c", filePath, "config file")
	pflag.Parse()
	viper.SetConfigType("yaml")
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
	conf := viper.GetString("config")
	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("load config %s fail: %v", conf, err))
	}

	// 解析初始配置
	baseConf := &Config{}
	if err := viper.Unmarshal(baseConf); err != nil {
		if err != nil {
			panic(err)
		}
	}

	// AuthExpireTime 解析为 time.Duration
	authExpireTime, err := time.ParseDuration(baseConf.Auth.ExpireTime)
	if err != nil {
		panic(err)
	}

	// 构造 SugaredConfig
	sugaredConfig := &SugaredConfig{
		Config:         baseConf,
		AuthExpireTime: authExpireTime,
	}

	return sugaredConfig
}

// 讀取環境變數 .env 或 docker-compose 的
func NewConfigFromEnv() *SugaredConfig {

	// 讀取環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	maxSize, err := utils.ConvertStringToInt(os.Getenv("LOG_MAX_SIZE"))
	if err != nil {
		log.Panicf("error=%v", err)
	}
	maxAge, err := utils.ConvertStringToInt(os.Getenv("LOG_MAX_AGE"))
	if err != nil {
		log.Panicf("error=%v", err)
	}
	maxBackups, err := utils.ConvertStringToInt(os.Getenv("LOG_BACKUPS"))
	if err != nil {
		log.Panicf("error=%v", err)
	}

	baseConf := &Config{
		Web: Web{
			Mode: os.Getenv("WEB_MODE"),
			Port: os.Getenv("WEB_PORT"),
		},
		Mysql: Mysql{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Auth: Auth{
			Active:     os.Getenv("AUTH_ACTIVE"),
			ExpireTime: os.Getenv("AUTH_EXPIRE_TIME"),
			PrivateKey: os.Getenv("AUTH_PRIVATE_KEY"),
		},
		Redis: Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Log: Log{
			Env:        os.Getenv("WEB_MODE"),
			Path:       os.Getenv("LOG_PATH"),
			Encoding:   os.Getenv("LOG_ENCODING"),
			MaxSize:    maxSize,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
		},
	}

	//AuthExpireTime 解析为 time.Duration
	authExpireTime, err := time.ParseDuration(baseConf.Auth.ExpireTime)
	if err != nil {
		panic(err)
	}

	// 构造 SugaredConfig
	sugaredConfig := &SugaredConfig{
		Config:         baseConf,
		AuthExpireTime: authExpireTime,
	}

	return sugaredConfig
}
