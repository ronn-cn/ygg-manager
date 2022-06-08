package main

import (
	"fmt"
	"libs/logger"
	N "libs/network"
	"libs/ouid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleSystem(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-system-list":
		getSystemList(c)
	case "query-system":
		querySystem(c)
	case "create-system":
		createSystem(c)
	case "update-system":
		updateSystem(c)
	case "delete-system":
		deleteSystem(c)
	case "push-system":
		pushSystem(c)
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
	if _, err := VerifyToken(c); err == nil {
		namestr := c.Query("name")
		var systems []System

		var tx *gorm.DB = PGDB
		if namestr != "" {
			tx = PGDB.Where("name LIKE ?", "%"+namestr+"%")
		}
		if result := tx.Order("id").Find(&systems); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(systems), "items": systems}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 查询系统
func querySystem(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if _, err := VerifyToken(c); err == nil {
		ouidstr := c.Query("ouid")

		if ouidstr == "" {
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
			return
		}

		var system System
		if result := PGDB.Where("ouid = ?", ouidstr).Find(&system); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": system})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
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
			if system.OUID == "" {
				system.OUID = ouid.GenerateOUID()
			}
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

// 更新系统
func updateSystem(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		var system System
		if err := c.BindJSON(&system); err == nil {
			if result := PGDB.Debug().Model(&System{OUID: system.OUID}).Select("name", "list", "main", "remark").Updates(&system); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": system})
			} else {
				c.JSON(200, gin.H{"errcode": 10204, "errmsg": "更新数据错误"})
			}
		} else {
			logger.Debugf("%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 删除系统
func deleteSystem(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		ouidstr := c.Query("ouid")
		logger.Debugf("请求的OUID：%v", ouidstr)
		if result := PGDB.Debug().Where("ouid = ?", ouidstr).Delete(&System{}); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": ouidstr})
		} else {
			fmt.Println("数据库返回的结构:", result)
			c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 推送系统
func pushSystem(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if cert, err := VerifyToken(c); err == nil {
		fmt.Println(cert)
		var system System
		if err := c.BindJSON(&system); err == nil {
			PGDB.Debug().Where("ouid = ?", system.OUID).First(&system)
			var devices []Device
			if result := PGDB.Debug().Where("system = ?", system.OUID).Find(&devices); result.Error == nil {
				pushdata := make(map[string]interface{})
				deviceouids := make([]string, 0)
				for _, dev := range devices {
					deviceouids = append(deviceouids, dev.OUID)
				}
				pushdata["devices"] = deviceouids
				pushdata["system"] = system

				obj := N.NewOMTPDataToJSON("PushUpdate", N.REPORT, 0, pushdata)
				if msgBytes, err := N.NewOMTPMessage("json", obj); err == nil {
					if err := SelfNode.NodeSend(GatherNode, msgBytes); err == nil {
						logger.Infof("[%s]: %v", SelfNode.OUID, "推送："+string(msgBytes))
						c.JSON(200, gin.H{"errcode": 0, "errmsg": "推送更新成功"})
						return
					} else {
						logger.Errorf("[%s]: %v", SelfNode.OUID, err)
					}
				} else {
					logger.Errorf("[%s]: %v", "推送", err)
				}
			}

			c.JSON(200, gin.H{"errcode": 10205, "errmsg": "推送更新失败"})
			return
		}
	}
}
