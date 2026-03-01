package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yuhang-jieke/consulls/srv/getaway/handler"
)

func Router(orderHandler *handler.OrderHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/")
	{
		api.POST("/orders", orderHandler.GoodsAdd)
		api.POST("/update", orderHandler.GoodsUpdate)
		api.PUT("/orders/:id")
		api.DELETE("/orders/:id")
	}

	return r
}
