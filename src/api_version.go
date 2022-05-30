package main

import (
	"libs/logger"

	"github.com/gin-gonic/gin"
)

func handleVersion(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-version-list":
		// getApplicationList(c)
	case "query-version":
		// queryApplication(c)
	case "create-version":
		createVersion(c) // 新增版本
	default:
		c.Status(404)
		return
	}
}

func queryVersion(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
}

// 创建版本
func createVersion(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}
	if _, err := VerifyToken(c); err == nil {
		var version Version
		if err := c.BindJSON(&version); err == nil {
			// version.CreatorID = &cert.OUID
			if result := PGDB.Debug().Create(&version); result.Error == nil {
				// 创建新版本时，将应用表中的最新版本更新
				PGDB.Debug().Model(&Application{}).Where("appid = ?", version.Appid).Update("latest", version.Version)
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": version.Version})
			} else {
				c.JSON(200, gin.H{"errcode": 10203, "errmsg": "创建数据错误"})
			}
		} else {
			logger.Debugf("%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}
