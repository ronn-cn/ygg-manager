package main

import (
	"fmt"
	"libs/logger"

	"github.com/gin-gonic/gin"
)

func handleCompany(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-company-list": // 查询商户列表
		getCompanyList(c)
	case "create-company":
		createCompany(c)
	case "update-company":
		updateCompany(c)
	case "delete-company":
		deleteCompany(c)
	default:
		c.Status(404)
		return
	}
}

func getCompanyList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		var companies []Company // 公司复数
		if result := PGDB.Order("id").Find(&companies); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(companies), "items": companies}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

func createCompany(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)

		var company Company
		if err := c.BindJSON(&company); err == nil {
			if result := PGDB.Debug().Create(&company); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": company.Name})
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

func updateCompany(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)

		var company Company
		if err := c.BindJSON(&company); err == nil {
			if result := PGDB.Debug().Model(&Company{ID: company.ID}).Updates(&company); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": company})
			} else {
				c.JSON(200, gin.H{"errcode": 10204, "errmsg": "更新数据错误"})
			}
		} else {
			logger.Debugf("%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

// 删除商户
func deleteCompany(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)

		idstr := c.Query("id")
		logger.Debugf("请求的ID：%v", idstr)
		if result := PGDB.Debug().Where("id = ?", idstr).Delete(&Company{}); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": idstr})
		} else {
			fmt.Println("数据库返回的结构:", result)
			c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}
