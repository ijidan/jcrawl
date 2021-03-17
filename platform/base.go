package platform

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/ijidan/jcrawl/model"
	"gorm.io/gorm"
	"strings"
)

//基础类
type BasePlatform struct {
	C  *colly.Collector
	Db *gorm.DB
}


//数据处理
func (b *BasePlatform) handleRecord(db *gorm.DB, housePlatform string, houseId string, houseName string, houseInfoStr string, houseUrl string, updateTime string) {
	//查询
	var summary model.HouseSummary
	result := db.Where("house_platform= ? and house_id=?", housePlatform, houseId).First(&summary)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		//插入数据
		summaryModel := model.HouseSummary{
			HousePlatform: housePlatform,
			HouseId:       houseId,
			HouseName:     houseName,
			HouseInfo:     houseInfoStr,
			HouseUrl:      houseUrl,
			UpdateAt:      updateTime,
		}
		db.Create(&summaryModel)
	} else {
		//更新数据
		summary.HousePlatform = housePlatform
		summary.HouseId = houseId
		summary.HouseName = houseName
		summary.HouseInfo = houseInfoStr
		summary.HouseUrl = houseUrl
		summary.UpdateAt = updateTime
		db.Save(&summary)
	}
}

//清理字符串
func (b *BasePlatform) cleanStr(str string) string {
	//去除HTML
	str = strip.StripTags(str)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	//替换引号
	str = strings.Replace(str, "：", "", -1)
	str = strings.Replace(str, ":", "", -1)
	//去除空格
	str = strings.TrimSpace(str)
	return str
}

//获取实例
func NewBasePlatform() *BasePlatform {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	ins := &BasePlatform{
		C:  c,
		Db: model.NewDb(),
	}
	return ins
}
