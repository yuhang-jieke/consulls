package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuhang-jieke/consulls/srv/getaway/client"
	"github.com/yuhang-jieke/consulls/srv/getaway/handler/request"
	__ "github.com/yuhang-jieke/consulls/srv/user-server/handler/proto"
)

type OrderHandler struct {
	orderClient *client.OrderClient
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderClient *client.OrderClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}
func (o *OrderHandler) GoodsAdd(c *gin.Context) {
	var form request.GoodsAdd
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "参数不正确",
		})
		return
	}

	grpcClient, err := o.orderClient.GetClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务不可用: " + err.Error(),
		})
		return
	}

	ctx := context.Background()
	_, err = grpcClient.AddGoods(ctx, &__.AddGoodsReq{
		Name:  form.Name,
		Price: form.Price,
		Stock: int64(form.Stock),
	})
	if err != nil {
		log.Printf("AddGoods failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "调用服务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "商品添加成功",
	})
}
func (o *OrderHandler) GoodsUpdate(c *gin.Context) {
	var form request.GoodsUPdate
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "参数不正确",
		})
		return
	}
	grpcClient, err := o.orderClient.GetClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务不可用: " + err.Error(),
		})
		return
	}
	_, err = grpcClient.UpdateGoods(c, &__.UpdateGoodsReq{
		Id:    int64(form.Id),
		Stock: int64(form.Stock),
	})
	if err != nil {
		log.Printf("AddGoods failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "调用服务失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "商品修改成功",
	})
}
