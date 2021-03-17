package model

import (
	"fmt"
	"github.com/ijidan/jcrawl/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//获取实例
func NewDb() *gorm.DB {
	conf:=pkg.NewConfig()
	hostname:=conf.Get("database.hostname")
	port:=conf.Get("database.port")
	username:=conf.Get("database.username")
	password:=conf.Get("database.password")
	database:=conf.Get("database.database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		panic("failed to connect database")
	}
	return db
}
