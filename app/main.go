package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)
import "github.com/joiller/url-shortener"

func main() {
	route := gin.Default()
	route.Use(gin.Logger(), gin.Recovery())
	group := route.Group("shorten")
	group.POST("/", url_shortener.LongToShortHandler)
	group.GET("/:short", url_shortener.ShortToLongHandler)
	group.PUT("/:short", url_shortener.UpdateShortUrlHandler)
	group.DELETE("/:short", url_shortener.DeleteShortUrlHandler)
	group.GET("/:short/stats", url_shortener.GetShortUrlStatusHandler)
	err := route.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
