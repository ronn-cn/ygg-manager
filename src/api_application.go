package main

import (
	"fmt"
	"libs/convert"
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
	case "query-application":
		queryApplication(c)
	case "query-application-version":
		queryApplicationVersion(c)
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

	if cert, err := VerifyToken(c); err == nil {
		fmt.Println(cert)

		pagestr := c.Query("page")
		limitstr := c.Query("limit")

		if pagestr != "" {
			pagestr = "1"
		}
		if limitstr != "" {
			limitstr = "9999"
		}

		namestr := c.Query("name")
		appidstr := c.Query("appid")
		typestr := c.Query("type")

		var applications []Application

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
		if result := tx.Debug().Order("id").Find(&applications); result.Error == nil {

			items := applications
			total := int64(len(applications))

			if pagestr != "" {
				limitnum := convert.StoI(limitstr)
				startnum := (convert.StoI(pagestr) - 1) * limitnum
				endnum := convert.StoI(pagestr) * limitnum
				if endnum > total {
					endnum = total
				}
				items = applications[startnum:endnum]
			}
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": total, "items": items}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		logger.Debugf("发现错误：%v", err)
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

// 查询应用
func queryApplication(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	appidstr := c.Query("appid")
	// 只要验证信息通过就能查询
	if _, err := VerifyToken(c); err == nil {
		var app Application
		if result := PGDB.Debug().Where("appid = ?", appidstr).First(&app); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": app})
			return
		} else {
			c.JSON(200, gin.H{"errcode": 10201, "errmsg": "数据错误"})
			return
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
		return
	}
}

func queryApplicationVersion(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	appidstr := c.Query("appid")
	if _, err := VerifyToken(c); err == nil {
		var versions []Version
		if result := PGDB.Debug().Order("id").Where("appid = ?", appidstr).Find(&versions); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(versions), "items": versions}})
			return
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
			return
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
		return
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
			if app.Appid == "" {
				app.Appid = MD5(ouid.GenerateOUID())
			}
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
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
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
