package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(idChan chan int) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "service run")
	})

	r.GET("/create", func(c *gin.Context) {
		// id写入通道
		select {
		case idChan <- 1:
			// 成功
			c.String(http.StatusOK, "[队列][写入]成功: 1")
		default:
			// Channel 已满，记录日志，走定时Task
			c.String(http.StatusOK, "[队列][写入]失败,idChan已满: 1,走定时Task")
		}
	})

	return r
}
