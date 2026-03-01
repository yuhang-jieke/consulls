package model

import (
	__ "github.com/yuhang-jieke/consulls/srv/user-server/handler/proto"
	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(30);comment:商品名称"`
	Price float64 `gorm:"type:decimal(10,2);comment:商品价格"`
	Stock int     `gorm:"type:int(10);comment:库存"`
}

func (g *Goods) GoodsAdd(db *gorm.DB) error {
	return db.Create(&g).Error
}

func (g *Goods) UpdateGoods(db *gorm.DB, in *__.UpdateGoodsReq) error {
	return db.Model(&Goods{}).Where("id=?", in.Id).Update("stock", in.Stock).Error
}
func (g *Goods) ListGoods(db *gorm.DB, in *__.SearchGoodsReq) ([]*Goods, error) {
	var list []*Goods
	if in.Size < 0 || in.Size > 3 {
		in.Size = 1
	}
	if in.Page <= 0 || in.Page > 3 {
		in.Page = 1
	}
	offset := (in.Page - 1) * in.Size
	err := db.Model(&Goods{}).Offset(int(offset)).Limit(int(in.Size)).Find(&list).Error
	return list, err
}

func (g *Goods) GetGoodsId(db *gorm.DB, id int64) error {
	return db.Model(&Goods{}).Where("id=?", id).First(&g).Error
}
