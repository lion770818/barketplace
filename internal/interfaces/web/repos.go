package web

import (
	"marketplace/internal/config"

	//"ddd_demo/internal/bill"
	//"ddd_demo/internal/user"

	"marketplace/internal/common/mysql"

	//  mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Repos struct {
	// UserRepo user.UserRepo
	// AuthRepo user.AuthInterface
	// BillRepo bill.BillRepo
}

func NewRepos(cfg *config.SugaredConfig) *Repos {
	// 持久化类型的 repo
	mysqlCfg := mysql.Config{
		LogMode:  cfg.Mysql.LogMode,
		Driver:   cfg.Mysql.Driver,
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		Database: cfg.Mysql.Database,
		User:     cfg.Mysql.User,
		Password: cfg.Mysql.Password,
	}
	db := mysql.NewDB(mysqlCfg)
	if db == nil {
		return nil
	}

	//userRepo := user.NewMysqlUserRepo(db)
	//billRepo := bill.NewMysqlBillRepo(db)

	// auth 策略
	// var authRepo user.AuthInterface
	// if cfg.Auth.Active == "redis" {
	// 	authRepo = user.NewRedisAuthRepo(redis.NewClient(cfg), cfg.AuthExpireTime)
	// } else {
	// 	authRepo = user.NewJwtAuth(cfg.Auth.PrivateKey, cfg.AuthExpireTime)
	// }

	return &Repos{
		// UserRepo: userRepo,
		// AuthRepo: authRepo,
		// BillRepo: billRepo,
	}
}
