package repositories

import (
	"marketplace/internal/common/mysql"
	"marketplace/internal/config"
	"marketplace/internal/domain/entity"
	"marketplace/internal/domain/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type RepositoriesManager struct {
	User repository.UserRepository
	Food repository.FoodRepository
	db   *gorm.DB
}

// 建立db連線
func NewRepositories(cfg *config.SugaredConfig) (*RepositoriesManager, error) {

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

	return &RepositoriesManager{
		User: NewUserRepository(db),
		Food: NewFoodRepository(db),
		db:   db,
	}, nil
}

// closes the  database connection
func (s *RepositoriesManager) Close() error {
	return s.db.Close()
}

func (s *RepositoriesManager) GetDB() *gorm.DB {
	return s.db
}

// This migrate all tables
func (s *RepositoriesManager) Automigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
}
