package model

import "gorm.io/gorm"

const HousePlatformITHaiHome = "1"
const HousePlatformHipFlat = "2"
const HousePlatformPlus = "3"
const HousePlatformDDProperty = "4"

//房源概况
type HouseSummary struct {
	gorm.Model
	HousePlatform string
	HouseId       string
	HouseName     string
	HouseInfo     string
	HouseUrl      string
	UpdateAt      string
}

//表名
func (h *HouseSummary) TableName() string {
	return "t_house_summary"
}

//房源详细
type HouseDetail struct {
	HousePlatform string
	HouseName     string
	HouseType     string
	HouseTrans    string
	HouseArea     string
	HouseDeco     string
	HouseFloor    string
	HouseAddress  string
	UpdateAt      string
}
