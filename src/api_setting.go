package main

import (
	"libs/logger"

	"github.com/gin-gonic/gin"
)

func handleSetting(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-setting": // 查询应用列表
		getSetting(c)
	case "save-setting": //保存设置
		saveSetting(c)
	default:
		c.Status(404)
		return
	}
}

func getSetting(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	// 查询账户权限
	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		if account.TypeID == 1 {
			var settings []Setting // 设置复数
			PGDB.Find(&settings)
			redata := make(map[string]interface{})
			for _, item := range settings {
				redata[item.Key] = item.Value
			}
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": redata})
		} else {
			// 没有权限访问 10104
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问", "data": nil})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

func saveSetting(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	// 查询账户权限
	if account, err := VerifyToken(c); err == nil {
		logger.Debugf("请求的账号信息:%v", account)
		if account.TypeID == 1 {
			if err := c.BindJSON(&Config); err == nil {
				configToDB()   // 配置写入数据库
				configToFile() // 配置写入配置文件
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功"})
			} else {
				logger.Debugf("%v", err)
				c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
			}
		} else {
			// 没有权限访问 10104
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问", "data": nil})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}
