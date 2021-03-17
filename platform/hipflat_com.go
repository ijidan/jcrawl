package platform

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gocolly/colly/v2"
	"github.com/ijidan/jcrawl/model"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//结构体
type HipFlat struct {
	*BasePlatform
}

//入口函数
func (p *HipFlat) Start() error {
	perPage := 34
	c := p.C
	var appId string
	var apiKey string
	var indexPrefix string
	// 异常处理
	c.OnError(func(response *colly.Response, err error) {
		log.Println(err.Error())
	})
	//请求之前
	c.OnRequest(func(request *colly.Request) {
		log.Println("开始抓取入口页： ", request.URL.String())
	})
	c.OnResponse(func(response *colly.Response) {
		body := response.Body
		content := string(body)
		if appId == "" {
			appId, apiKey, indexPrefix = p.ParseToken(content)
		}
		for page := 1; page <= 10000; page++ {
			cnt := p.ParseList(c, apiKey, appId, page, perPage)
			if cnt < perPage {
				break
			}
		}
		log.Println("抓取结束")
	})
	startUrl := "https://www.hipflat.com/search/sale/condo,house,townhouse_y/TH.BM_r1/any_r2/any_p/any_b/any_a/any_w/any_i/100.6244261045141,13.77183154691727_c/12_z/list_v"
	return p.C.Visit(startUrl)
}

//解析TOKEN
func (p *HipFlat) ParseToken(content string) (string, string, string) {
	pattern := regexp.MustCompile(`appId.*},`)
	matched := pattern.FindStringSubmatch(content)
	if len(matched) == 0 {
		return "", "", ""
	}
	idStr := matched[0]
	idStr = `{"` + strings.Trim(idStr, ",")
	var tempMap map[string]string
	err := json.Unmarshal([]byte(idStr), &tempMap)
	if err != nil {
		return "", "", ""
	}
	return tempMap["appId"], tempMap["apiKey"], tempMap["indexPrefix"]
}

//解析列表
func (p *HipFlat) ParseList(c *colly.Collector, apiKey string, appId string, page int, perPage int) int {
	cnt := 0
	url := "https://tom76npq59-dsn.algolia.net/1/indexes/PROD_ListingSearchResults/query?x-algolia-api-key=" + apiKey + "&x-algolia-application-id=" + appId + "&x-algolia-agent=Algolia%20for%20jQuery%203.10.2"
	collector := c.Clone()
	//请求之前
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})
	//请求结束
	collector.OnResponse(func(response *colly.Response) {
		body := response.Body
		content := string(body)
		cnt = p.ParseDetail(content)
		logContent := fmt.Sprintf("抓取条数：%d,抓取URL：%s", cnt, response.Request.URL.String())
		log.Print(logContent)
	})
	//错误
	collector.OnError(func(response *colly.Response, err error) {
		logContent := fmt.Sprintf("抓取错误：%s，抓取URL：%s", err.Error(), response.Request.URL.String())
		log.Println(logContent)
	})
	bodyContent := `{"params":"query=&facets=*&facetFilters=transaction%3Asale%2C(property_type%3Acondo%2Cproperty_type%3Ahouse%2Cproperty_type%3Atownhouse)%2Cregion.hasc%3ATH.BM&numericFilters=&page=` + strconv.Itoa(page) + `&hitsPerPage=` + strconv.Itoa(perPage) + `"}`
	payload := []byte(bodyContent)
	_ = collector.PostRaw(url, payload)
	return cnt
}

//解析详情
func (p *HipFlat) ParseDetail(content string) int {
	js, _ := simplejson.NewJson([]byte(content))
	dataList, _ := js.Get("hits").Array()
	dataCnt := len(dataList)
	for _, v := range dataList {
		vMap := v.(map[string]interface{})
		//slug_en
		slugEn := vMap["slug"].(string)
		//项目
		//projectName:=vMap["project_name"]
		houseId := vMap["objectID"].(string)
		title := vMap["title"]
		titleMap := title.(map[string]interface{})
		titleEn := titleMap["en"].(interface{})
		houseName := titleEn.(string)
		houseInfoJson, _ := json.Marshal(v)
		houseInfoStr := string(houseInfoJson)

		houseUrl := fmt.Sprintf("https://www.hipflat.com/listings/%s", slugEn)
		updateTime := ""
		//处理数据
		p.BasePlatform.handleRecord(p.Db, model.HousePlatformHipFlat, houseId, houseName, houseInfoStr, houseUrl, updateTime)
	}
	return dataCnt
}

//获取实例
func NewHipFlat() *HipFlat {
	ins := &HipFlat{
		BasePlatform: NewBasePlatform(),
	}
	return ins
}
