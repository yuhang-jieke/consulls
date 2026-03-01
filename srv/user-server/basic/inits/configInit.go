package inits

import (
	"github.com/spf13/viper"
	"github.com/yuhang-jieke/consulls/srv/user-server/basic/config"
)

func ConfigInit() {
	viper.SetConfigFile("C:\\Users\\ZhuanZ\\Desktop\\zuoye7\\consulls\\srv\\dev.yaml")
	viper.ReadInConfig()
	viper.Unmarshal(&config.GlobalConf)
}
