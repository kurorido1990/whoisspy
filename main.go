package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"reflect"

	//"net/http"
)

const (
	Status_OK = 200
)

var topic = map[int][]string {
	1 : []string{
		"ipad", "iphone",
	},
	2 : []string{
		"包子", "餃子",
	},
	3 : []string{
		"康熙", "乾隆",
	},
	4 : []string{
		"天主教", "基督教",
	},
	5 : []string{
		"浴缸", "魚缸",
	},
	6 : []string{
		"電動車", "摩托車",
	},
	7 : []string{
		"眉毛", "睫毛",
	},
	8 : []string{
		"筷子", "竹籤",
	},
	9 : []string{
		"麻雀", "烏鴉",
	},
	10 : []string{
		"鏡子", "玻璃",
	},
	11 : []string{
		"那英", "王菲",
	},
	12 : []string{
		"樹枝", "樹幹",
	},
	13 : []string{
		"牛奶", "豆漿",
	},
	14 : []string{
		"香港", "台灣",
	},
	15 : []string{
		"辣椒", "芥末",
	},
	16 : []string{
		"海豚", "海獅",
	},
	17 : []string{
		"蝴蝶", "蜜蜂",
	},
	18 : []string{
		"首爾", "東京",
	},
	19 : []string{
		"柳丁", "橘子",
	},
	20 : []string{
		"新年", "跨年",
	},
	21 : []string{
		"吉他", "Bass",
	},
	22 : []string{
		"公車", "地鐵",
	},
	23 : []string{
		"火車", "高鐵",
	},
	24 : []string{
		"結婚", "訂婚",
	},
	25 : []string{
		"情人節", "光棍節",
	},
}

type IndexData struct {
	Title   string
	Content string
}

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一個首頁"
	c.HTML(http.StatusOK, "index.html", data)
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/*")
	server.GET("/", test)
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(Status_OK, gin.H{
			"message" : "pong",
		})
	})
	server.GET("/topic", func(context *gin.Context) {
		context.JSON(Status_OK, gin.H{
			"topic" : getTopic(),
		})
	})
	//server.Run(":9999")
	server.Run()
}

func getTopic() string {
	return topic[MapRandomKeyGet(topic).(int)][0]
}

func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}
