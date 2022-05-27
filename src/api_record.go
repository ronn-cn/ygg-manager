package main

import (
	"libs/logger"

	"github.com/gin-gonic/gin"
)

func handleRecord(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-record-list": // 查询应用列表
		getRecordList(c)
	case "create-record":
		createRecord(c)
	default:
		c.Status(404)
		return
	}
}

func getRecordList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	if _, err := VerifyToken(c); err == nil {
		var records []Record // 日志复数
		if result := PGDB.Debug().Find(&records); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(records), "items": records}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 创建记录
func createRecord(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if _, err := VerifyToken(c); err == nil {
		var record Record
		if err := c.BindJSON(&record); err == nil {
			if result := PGDB.Debug().Create(&record); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功"})
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
