package main

import (
	"fmt"
	"libs/logger"
	"libs/ouid"

	"github.com/gin-gonic/gin"
)

func handleLicense(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-license-list": // 查询密钥列表
		getLicenseList(c)
	case "create-license": // 创建密钥
		createLicense(c)
	case "delete-license":
		deleteLicense(c)
	default:
		c.Status(404)
		return
	}
}

// 查询密钥列表
func getLicenseList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		var licenses []License

		if result := PGDB.Order("id").Find(&licenses); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(licenses), "items": licenses}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 创建密钥
func createLicense(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		var license License
		if err := c.BindJSON(&license); err == nil {
			license.Code = MD5(ouid.GenerateOUID())
			license.CreatorID = &account.OUID
			if result := PGDB.Debug().Create(&license); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": license.Code})
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

func deleteLicense(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		codestr := c.Query("code")
		logger.Debugf("请求的Code：%v", codestr)
		if result := PGDB.Debug().Where("code = ?", codestr).Delete(&License{}); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": codestr})
		} else {
			fmt.Println("数据库返回的结构:", result)
			c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}
