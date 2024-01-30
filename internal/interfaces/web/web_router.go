package web

import (
	"marketplace/internal/common/logs"
	"marketplace/internal/interfaces/middleware"
	//"bito_group/internal/user"
)

func WithRouter(s *WebServer) {

	logs.Debugf("註冊路由")

	// 新建 handler
	userHandler := s.Apps.NewUserHandler()
	foodHandler := s.Apps.NewFoodHandler()
	authHandler := s.Apps.NewAuthHandler()

	// api
	api := s.Engin.Group("/v1")

	// 中间件
	//api.Use(authMiddleware.Auth)

	// 路由
	api.POST("/users", userHandler.SaveUser)
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:user_id", userHandler.GetUser)

	//post routes
	api.POST("/food", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), foodHandler.SaveFood)
	api.PUT("/food/:food_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), foodHandler.UpdateFood)
	api.GET("/food/:food_id", foodHandler.GetFoodAndCreator)
	api.DELETE("/food/:food_id", middleware.AuthMiddleware(), foodHandler.DeleteFood)
	api.GET("/food", foodHandler.GetAllFood)

	// authentication routes
	api.POST("/login", authHandler.Login)
	api.POST("/logout", authHandler.Logout)
	api.POST("/refresh", authHandler.Refresh)
}
