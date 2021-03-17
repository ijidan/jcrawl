package platform

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gocolly/colly/v2"
	"github.com/ijidan/jcrawl/model"
	"log"
	"strconv"
	"time"
)

//结构体
type Plus struct {
	*BasePlatform
}

//入口函数
func (p *Plus) Start() error {
	perPage:=24
	c := p.C
	c.OnRequest(func(request *colly.Request) {
		log.Print("start....")
	})
	c.OnResponse(func(response *colly.Response) {
		for page := 1; page <= 10000; page++ {
			cnt:= p.ParseList(c, page,perPage)
			if cnt<perPage{
				break
			}
		}
		log.Print("抓取结束")
	})
	c.OnError(func(response *colly.Response, err error) {
		log.Print("抓取错误"+err.Error())
	})
	startUrl := "https://www.plus.co.th/search?lang=en&od_type=sale&q=&stock_type_id=&min_selling_price=&max_selling_price=&min_rental_price=&max_rental_price=&page=1&per_page=24&t=1614754479&map=false"
	return c.Visit(startUrl)
}

//解析列表
func (p *Plus) ParseList(c *colly.Collector, page int,perPage int) int {
	url := "https://agency-api.plus.co.th/api/v1/plus/search"
	collector := c.Clone()
	cnt:=0
	//请求之前
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Host", "agency-api.plus.co.th")
		request.Headers.Set("Accept", "application/json, text/plain, */*")
		request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36")
		request.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Headers.Set("Origin", "https://www.plus.co.th")
		request.Headers.Set("Referer", "ttps://www.plus.co.th/")
	})
	//请求结束
	collector.OnResponse(func(response *colly.Response) {
		body := response.Body
		content := string(body)
		cnt=p.ParseDetail(content)
		log.Print("列表抓取条数："+strconv.Itoa(cnt))
	})
	//错误
	collector.OnError(func(response *colly.Response, err error) {
		log.Println("列表抓取错误："+err.Error())
	})
	bodyContent := fmt.Sprintf("lang=en&od_type=sale&page=%d&per_page=%d&t=%d&map=false", page,perPage,time.Now().Second())
	payload := []byte(bodyContent)
	_=collector.PostRaw(url, payload)
	return  cnt
}

//详情
func (p *Plus) ParseDetail(content string) int {
	js, _ := simplejson.NewJson([]byte(content))
	dataList,_:=js.Get("data").Get("data").Array()
	dataCnt:=len(dataList)
	for _,v:=range dataList{
		vMap:=v.(map[string]interface{})
		//houseId:=vMap["id"].(string)
		//属性
		attr:=vMap["attributes"]
		attrMap:=attr.(map[string]interface{})

		//slug_en
		headlineEn:=attrMap["stock_headline_en"].(string)
		slugEn:=attrMap["slug_en"].(string)
		plusRef:=attrMap["plus_ref"].(string)
		//项目
		//project:=attrMap["project"]
		//projectMap:=project.(map[string]interface{})
		houseId:=slugEn
		houseName:=headlineEn
		houseInfoJson,_ :=json.Marshal(v)
		houseInfoStr :=string(houseInfoJson)

		houseUrl:=fmt.Sprintf("https://www.plus.co.th/unit/%s/%s?lang=en",slugEn,plusRef)
		updateTime:=""

		//处理数据
		p.BasePlatform.handleRecord(p.Db, model.HousePlatformPlus, houseId, houseName, houseInfoStr, houseUrl, updateTime)
	}
	return dataCnt
}

//获取实例
func NewPlus() *Plus {
	ins := &Plus{
		BasePlatform: NewBasePlatform(),

	}
	return ins
}
