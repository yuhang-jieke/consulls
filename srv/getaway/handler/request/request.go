package request

type GoodsAdd struct {
	Name  string  `form:"name"   binding:"required"`
	Price float64 `form:"price" binding:"required"`
	Stock int     `form:"stock" binding:"required"`
}
type GoodsUPdate struct {
	Id    int `form:"id" binding:"required"`
	Stock int `form:"stock" binding:"required"`
}
