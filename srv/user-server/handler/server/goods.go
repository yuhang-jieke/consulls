package server

import (
	"context"
	"errors"
	"log"

	"github.com/yuhang-jieke/consulls/srv/user-server/basic/config"
	__ "github.com/yuhang-jieke/consulls/srv/user-server/handler/proto"
	"github.com/yuhang-jieke/consulls/srv/user-server/model"
)

type Server struct {
	__.UnimplementedEcommerceServiceServer
}

func (s *Server) AddGoods(_ context.Context, in *__.AddGoodsReq) (*__.AddGoodsResp, error) {
	log.Printf("Received: %v", in.GetName())
	goods := model.Goods{
		Name:  in.Name,
		Price: in.Price,
		Stock: int(in.Stock),
	}
	err := (&goods).GoodsAdd(config.DB)
	if err != nil {
		return nil, errors.New("商品添加失败")
	}
	return &__.AddGoodsResp{
		Message: "商品添加成功",
	}, nil
}

func (s *Server) UpdateGoods(_ context.Context, in *__.UpdateGoodsReq) (*__.UpdateGoodsResp, error) {
	var goods model.Goods
	err := goods.UpdateGoods(config.DB, in)
	if err != nil {
		return nil, errors.New("商品修改失败")
	}
	return &__.UpdateGoodsResp{
		Message: "商品修改成功",
	}, nil
}
func (s *Server) SearchGoods(_ context.Context, in *__.SearchGoodsReq) (*__.SearchGoodsResp, error) {
	var goods model.Goods
	var listgoods []*model.Goods
	listgoods, err := goods.ListGoods(config.DB, in)
	if err != nil {
		return nil, errors.New("列表展示失败")
	}
	var list []*__.Goods
	for _, listgood := range listgoods {
		list = append(list, &__.Goods{
			Name:  listgood.Name,
			Price: listgood.Price,
			Stock: int64(listgood.Stock),
		})
	}
	return &__.SearchGoodsResp{
		Goods: list,
	}, nil
}
func (s *Server) GetGoodsById(_ context.Context, in *__.GetGoodsByIdReq) (*__.GetGoodsByIdResp, error) {
	var goods model.Goods
	err := goods.GetGoodsId(config.DB, in.Id)
	if err != nil {
		return nil, errors.New("获取失败")
	}
	return &__.GetGoodsByIdResp{
		Goods: &__.Goods{
			Name:  goods.Name,
			Price: goods.Price,
			Stock: int64(goods.Stock),
		},
	}, nil
}
