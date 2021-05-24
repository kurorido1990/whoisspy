package main

import (
	"github.com/gin-gonic/gin"
	//"net/http"
)

type IndexData struct {
	Title   string
	Content string
}

//func test(c *gin.Context) {
//	data := new(IndexData)
//	data.Title = "首頁"
//	data.Content = "我的第一個首頁"
//	c.HTML(http.StatusOK, "index.html", data)
//}

func main() {
	server := gin.Default()
	//server.LoadHTMLGlob("template/*")
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "pong",
		})
	})
	server.Run(":8888")
}
