package inits

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/yuhang-jieke/consulls/srv/user-server/basic/config"
)

var Ctx = context.Background()

func RedisInit() {
	conf := config.GlobalConf.Redis
	Addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: conf.Password, // no password set
		DB:       conf.Database, // use default DB
	})

	err := config.Rdb.Ping(Ctx).Err()
	if err != nil {
		panic("redis连接失败")
	}
	fmt.Println("redis连接成功")
}
