package main

import (
	"libs/convert"
	"libs/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

		pagestr := c.Query("page")
		limitstr := c.Query("limit")
		typestr := c.Query("type")
		levelstr := c.Query("level")
		infostr := c.Query("info")

		var tx *gorm.DB = PGDB
		if typestr != "" {
			tx = tx.Where("type = ?", typestr)
		}
		if infostr != "" {
			tx = tx.Where("info LIKE ?", "%"+infostr+"%")
		}
		if levelstr != "" {
			tx = tx.Where("level = ?", levelstr)
		}

		var records []Record // 日志复数
		if result := tx.Debug().Find(&records); result.Error == nil {
			items := records
			total := int64(len(records))

			if pagestr != "" {
				limitnum := convert.StoI(limitstr)
				startnum := (convert.StoI(pagestr) - 1) * limitnum
				endnum := convert.StoI(pagestr) * limitnum
				if endnum > total {
					endnum = total
				}
				items = records[startnum:endnum]
			}
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": total, "items": items}})
		} else {
			c.JSON(200, gin.H{"errcode": 10200, "errmsg": "查询数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
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
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}
