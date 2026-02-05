package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Cron struct {
	scheduler gocron.Scheduler
}

func NewCron() *Cron {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	if err != nil {
		fmt.Printf("[定时器]创建失败: %v\n", err)
		return nil
	}
	fmt.Printf("[定时器]创建成功\n")

	return &Cron{
		scheduler: scheduler,
	}
}

func (cron *Cron) NewCronTask() {
	// Task
	_, err := cron.scheduler.NewJob(
		gocron.DurationJob(
			time.Second*10,
		),
		gocron.NewTask(
			func() {
				fmt.Printf("[Task]逻辑处理\n")
			},
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		fmt.Printf("[定时器]任务Task创建失败: %v\n", err)
		return
	}
	fmt.Printf("[定时器]任务Task创建成功\n")

	cron.scheduler.Start()
}

func (c *Cron) Stop() {
	// Shutdown 会优雅停止：等待正在运行的任务完成
	if err := c.scheduler.Shutdown(); err != nil {
		fmt.Printf("[定时器]关闭时出错: %v\n", err)
	}
}
