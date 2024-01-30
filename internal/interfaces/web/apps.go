package web

import (
	"marketplace/internal/infrastructure/auth"
	"marketplace/internal/infrastructure/repositories"
	"marketplace/internal/interfaces"
	"marketplace/internal/interfaces/fileupload"
)

type Apps struct {
	// 不同的 use case 可以掛載
	Repos        *repositories.RepositoriesManager
	RedisService *auth.RedisService

	Users        *interfaces.Users
	Food         *interfaces.Food
	Authenticate *interfaces.Authenticate
}

func NewApps(repos *repositories.RepositoriesManager, redisService *auth.RedisService) *Apps {

	tk := auth.NewToken()
	fd := fileupload.NewFileUpload()

	users := interfaces.NewUsers(repos.User, redisService.Auth, tk)
	foods := interfaces.NewFood(repos.Food, repos.User, fd, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(repos.User, redisService.Auth, tk)

	return &Apps{
		Repos:        repos,
		RedisService: redisService,
		Users:        users,
		Food:         foods,
		Authenticate: authenticate,
	}
}

func (s *Apps) NewUserHandler() interfaces.UserInterface {
	return s.Users
}

func (s *Apps) NewFoodHandler() interfaces.FoodInterface {
	return s.Food
}

func (s *Apps) NewAuthHandler() interfaces.AuthInterface {
	return s.Authenticate
}
