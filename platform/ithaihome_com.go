package platform

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/ijidan/jcrawl/model"
	"log"
	"regexp"
	"strings"
)

//结构体
type IThaiHome struct {
	*BasePlatform
}

//入口函数
func (p *IThaiHome) Start() error {
	c := p.C
	// 异常处理
	c.OnError(func(response *colly.Response, err error) {
		log.Println(err.Error())
	})
	//请求之前
	c.OnRequest(func(request *colly.Request) {
		log.Println("抓取： ", request.URL.String())
	})
	//列表
	c.OnHTML("#example1_next > a.next", func(htmlElement *colly.HTMLElement) {
		href, found := htmlElement.DOM.Attr("href")
		if found {
			absUrl := htmlElement.Request.AbsoluteURL(href)
			err := htmlElement.Request.Visit(absUrl)
			if err != nil {
				log.Print("列表抓取错误："+err.Error())
			}
		}

	})
	//详情页链接
	c.OnHTML(".esfylist", func(htmlElement *colly.HTMLElement) {
		htmlElement.DOM.Find("a.fz14").Each(func(i int, selection *goquery.Selection) {
			url, found := selection.Attr("href")
			if found {
				absUrl := htmlElement.Request.AbsoluteURL(url)
				p.ParseDetail(c, absUrl)
			}
		})
	})
	startUrl := "http://www.ithaihome.com/Home/Fang/index.html?cid=2&type_id=3&city=47662"
	return p.C.Visit(startUrl)
}

//解析详情页
func (p *IThaiHome) ParseDetail(c *colly.Collector, url string) {
	//克隆
	collector := c.Clone()
	//请求之前
	var houseId string
	var houseUrl string
	//请求之前
	collector.OnRequest(func(request *colly.Request) {
		houseUrl = request.URL.String()
		houseId = p.ParseId(houseUrl)
		log.Println("详情页： ", houseUrl)

	})
	collector.OnHTML("body", func(htmlElement *colly.HTMLElement) {
		//名称
		houseName := htmlElement.DOM.Find(".info-main-hd > h1").Text()
		houseName = p.cleanStr(houseName)
		//最近更新
		updateTime := htmlElement.DOM.Find(".olorg").Text()
		updateTime = p.cleanStr(updateTime)
		updateTimeArr := strings.Split(updateTime, "|")
		updateTime = p.ParseUpdateTime(updateTimeArr)

		//房源信息
		houseInfoHd := htmlElement.DOM.Find(".house-info-bd > li")
		houseInfoList := map[string]string{}
		houseInfoHd.Each(func(i int, selection *goquery.Selection) {
			//key
			key := selection.Find("span.fcg").First().Text()
			key = p.cleanStr(key)
			//value
			value := selection.Find("i").Text()
			value = p.cleanStr(value)
			value = strings.Replace(value, " ", "", -1)
			if value == "" {
				value = selection.Find("span.fco").Text()
				value = p.cleanStr(value)
			}
			houseInfoList[key] = value
		})
		houseInfoJson, _ := json.Marshal(houseInfoList)
		houseInfoStr := string(houseInfoJson)

		//处理数据
		p.BasePlatform.handleRecord(p.Db, model.HousePlatformITHaiHome, houseId, houseName, houseInfoStr, houseUrl, updateTime)
	})
	err := collector.Visit(url)
	if err != nil {
		log.Println(err)
	}
}

//解析更新时间
func (p *IThaiHome) ParseUpdateTime(updateTimeAttr []string) string {
	updateTimeStr := updateTimeAttr[0]
	updateTimeStr = p.cleanStr(updateTimeStr)
	updateTimeStr = strings.Replace(updateTimeStr, "更新", "", -1)
	updateTimeStr = strings.Replace(updateTimeStr, "最近", "", -1)
	return updateTimeStr
}

//解析ID
func (p *IThaiHome) ParseId(url string) string {
	pattern := regexp.MustCompile(`(\d)+`)
	matched := pattern.FindStringSubmatch(url)
	if len(matched) == 0 {
		return "0"
	}
	idStr := matched[0]
	return idStr
}

//根据字段取值
func (p *IThaiHome) ParseValue(houseInfoList map[string]string, key string) string {
	value, ok := houseInfoList[key]
	if ok {
		return value
	}
	return ""
}

//获取实例
func NewIThaiHome() *IThaiHome {
	ins := &IThaiHome{
		BasePlatform: NewBasePlatform(),
	}
	return ins
}
