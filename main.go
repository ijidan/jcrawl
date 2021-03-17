package main

import (
	"flag"
	"fmt"
	"github.com/ijidan/jcrawl/platform"
	"log"
	"time"
)

//hipFlat爬取
func crawlHipFlat()  {
	flat:=platform.NewHipFlat()
	err := flat.Start()
	if err != nil {
		log.Print(err)
	}
}

//IThaiHome爬取
func crawlIThaiHome() {
	thai := platform.NewIThaiHome()
	err := thai.Start()
	if err != nil {
		log.Print(err)
	}
}

//Plus网站抓取
func crawlPlus() {
	plus := platform.NewPlus()
	err := plus.Start()
	if err != nil {
		log.Print(err)
	}
}


//入口函数
func main() {
	startTime := time.Now().Second()
	crawlHipFlat()
	flag.Parse()
	//p := flag.String("p", "", "平台")
	//switch *p {
	//case model.HousePlatformITHaiHome:
	//	crawlIThaiHome()
	//	break
	//case model.HousePlatformHipFlat:
	//	crawlPlus()
	//	break
	//case model.HousePlatformPlus:
	//	break
	//default:
	//	log.Fatalln("请输入平台参数")
	//	return
	//}
	endTime := time.Now().Second()
	timeDiff := endTime - startTime
	tip := fmt.Sprintf("耗时：%d秒", timeDiff)
	log.Println(tip)
}
