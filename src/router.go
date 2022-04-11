package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/view")
	})
	router.StaticFS("/view", http.Dir("./view"))
	router.Any("/manager-api/*path", handleManagerApi)
}
