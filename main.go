package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"snai.graceful.exit/cron"
	"snai.graceful.exit/router"
)

func init() {
	fmt.Printf("[启动]服务启动,监听端口: 8024\n")
	fmt.Printf("[启动]当前环境: 测试服\n")
}

type App struct {
	Task   *cron.Task
	Cron   *cron.Cron
	Router *gin.Engine
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	idChan := make(chan int, 100)
	app := &App{
		Task:   cron.NewTask(idChan),
		Cron:   cron.NewCron(),
		Router: router.NewRouter(idChan),
	}

	var wg sync.WaitGroup

	// 实时拉取链数据池,非阻塞
	app.Task.CreatePool(&wg)

	// 定时统计,非阻塞
	app.Cron.NewCronTask()

	// 3. 配置 HTTP 服务
	srv := &http.Server{
		Addr:    ":8024",
		Handler: app.Router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[启动]gin run failed: %v\n", err)
		}
	}()

	// --- 等待退出信号 ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("[退出]接收到退出信号，开始清理资源...\n")

	// 顺序 1: 先关 Gin（停止新数据进入）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("[退出]Gin 强制关闭: %v\n", err)
	}

	// 顺序 2: 停止 Cron（防止定时任务触发新的异步操作）
	app.Cron.Stop()
	fmt.Printf("[退出]Cron 已安全停止\n")

	// 顺序 C: 关闭 Channel 并等待 Pool（处理完存量数据）
	close(idChan)
	fmt.Printf("[退出]等待 Pool 处理剩余数据...\n")
	wg.Wait()

	fmt.Printf("[退出]所有组件已安全退出\n")
}
