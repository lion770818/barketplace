package signals

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Allenxuxu/gev/log"
)

// WaitSignal 等待信號, 收到信號後, 執行回調函式
func WaitWith(stop func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("優雅關閉服務...")

	// 优雅退出
	stop()
}
