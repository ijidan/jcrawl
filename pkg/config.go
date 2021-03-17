package pkg

import (
	"github.com/spf13/viper"
	"os"
)

//获取配置
func NewConfig() *viper.Viper {
	wd, _ := os.Getwd()
	pathSep := GetPathSep()
	confPath := wd + pathSep + "config"
	v := viper.New()
	v.AddConfigPath(confPath)
	v.SetConfigName("default")
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}

//获取文件分隔符
func GetPathSep() string {
	pathSep := string(os.PathSeparator)
	return pathSep
}
