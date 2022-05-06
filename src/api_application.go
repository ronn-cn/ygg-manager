package main

import (
	"libs/logger"
	"libs/ouid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleApplication(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-application-list":
		getApplicationList(c)
	case "create-application":
		createApplication(c)
	case "update-application":
		updateApplication(c)
	case "delete-application":
		deleteApplication(c)
	default:
		c.Status(404)
		return
	}
}

// 查询应用列表
func getApplicationList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		namestr := c.Query("name")
		appidstr := c.Query("appid")
		typestr := c.Query("type")

		var application []Application

		var tx *gorm.DB = PGDB
		if namestr != "" {
			tx = PGDB.Where("name LIKE ?", "%"+namestr+"%")
		}
		if appidstr != "" {
			tx = PGDB.Where("appid LIKE ?", "%"+appidstr+"%")
		}
		if typestr != "" {
			tx = PGDB.Where("type = ?", typestr)
		}
		if result := tx.Debug().Order("id").Find(&application); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(application), "items": application}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 创建应用
func createApplication(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		var app Application
		if err := c.BindJSON(&app); err == nil {
			app.Appid = MD5(ouid.GenerateOUID())
			app.CreatorID = &account.OUID
			if result := PGDB.Debug().Create(&app); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": app.Appid})
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

func updateApplication(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.Status(405)
		return
	}
	var app Application
	if err := c.BindJSON(&app); err == nil {
		if result := PGDB.Debug().Model(&Application{Appid: app.Appid}).Select("Name", "Type", "Latest", "Status", "Depend", "Remark").Updates(&app); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": app.Appid})
		} else {
			c.JSON(200, gin.H{"errcode": 10204, "errmsg": "更新数据错误"})
		}
	} else {
		logger.Debugf("绑定参数错误：%v", err)
		c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
	}
}

func deleteApplication(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}
	// 删除账号
	appidstr := c.Query("appid")
	logger.Debugf("请求的需要删除的APPID：%v", appidstr)
	if result := PGDB.Debug().Where("appid = ?", appidstr).Delete(&Application{}); result.Error == nil {
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": appidstr})
	} else {
		c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
	}
}
