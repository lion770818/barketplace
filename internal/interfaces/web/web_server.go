package web

import (
	"context"
	"fmt"
	"marketplace/internal/common/logs"
	"marketplace/internal/common/servers"
	"marketplace/internal/config"
	"net/http"
	"time"

	//_ "marketplace/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type WebServer struct {
	httpServer *http.Server
	Engin      *gin.Engine
	Apps       *Apps
}

func (s *WebServer) AsyncStart() {
	logs.Debugf("[服務啟動] [web] 服務地址: %s", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Fatalf("[服務啟動] [web] 服務異常: %s", zap.Error(err))
		}
	}()
}

func (s *WebServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logs.Debugf("[服務關閉] [web] 關閉服務")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logs.Fatalf("[服務關閉] [web] 關閉服務異常: %s", zap.Error(err))
	}

	s.CloseResource()
}

// 關閉資源
func (s *WebServer) CloseResource() {
	if err := s.Apps.Repos.GetDB().Close(); err != nil {
		logs.Fatalf("[服務關閉] [web] 關閉db異常: %s", zap.Error(err))
	}

	if err := s.Apps.RedisService.Client.Close(); err != nil {
		logs.Fatalf("[服務關閉] [web] 關閉redis異常: %s", zap.Error(err))
	}

	logs.Debugf("[服務關閉] [web] 關閉服務資源成功")
}

func NewWebServer(cfg *config.SugaredConfig, apps *Apps) servers.ServerInterface {

	logs.Debugf("創建 web server mode:%s poet:%s", cfg.Web.Mode, cfg.Web.Port)

	gin.SetMode(cfg.Web.Mode)
	e := gin.Default()
	e.Use(cors.Default())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Web.Port),
		Handler: e,
	}

	s := &WebServer{
		httpServer: httpServer,
		Engin:      e,
		Apps:       apps,
	}

	// 設定swgger
	urlStr := "http://localhost:" + cfg.Web.Port + "/swagger/doc.json"
	url := ginSwagger.URL(urlStr) // The url pointing to API definition
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 注册路由
	WithRouter(s)

	return s
}
