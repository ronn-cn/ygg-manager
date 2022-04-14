package main

import (
	"libs/logger"
	"libs/ouid"

	"github.com/gin-gonic/gin"
)

func handleSystem(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-system-list":
		getSystemList(c)
	case "create-system":
		createSystem(c)
	default:
		c.Status(404)
		return
	}
}

// 查询系统列表
func getSystemList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	var systems []System
	PGDB.Order("id").Find(&systems)
	c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(systems), "items": systems}})
}

func createSystem(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		var system System
		if err := c.BindJSON(&system); err == nil {
			system.OUID = ouid.GenerateOUID()
			system.CreatorID = &account.OUID
			if result := PGDB.Debug().Create(&system); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": system.OUID})
			} else {
				c.JSON(200, gin.H{"errcode": 10203, "errmsg": "创建数据错误"})
			}
		} else {
			logger.Debugf("绑定参数错误：%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}
