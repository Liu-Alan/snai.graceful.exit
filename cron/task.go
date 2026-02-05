package cron

import (
	"fmt"
	"sync"
)

type Task struct {
	idChan chan int
}

func NewTask(idChan chan int) *Task {
	return &Task{
		idChan: idChan,
	}
}

// 实时工作池,消费队列,主线程非阻塞
func (task *Task) CreatePool(wg *sync.WaitGroup) {
	// 多协程消费
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(workerID int, idChan chan int) {
			defer wg.Done()
			fmt.Printf("[协程-%v]已启动\n", i)
			// 遍历管道数据
			for id := range idChan {
				fmt.Printf("逻辑处理:%v\n", id)
			}
			fmt.Printf("[协程-%v]已安全退出\n", workerID)
		}(i, task.idChan)
	}
}
