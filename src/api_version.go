package main

import (
	"fmt"
	"libs/convert"
	"libs/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleVersion(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-version-list":
		getVersionList(c)
	case "query-version":
		queryVersion(c)
	case "create-version":
		createVersion(c) // 新增版本
	default:
		c.Status(404)
		return
	}
}

func getVersionList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	if cert, err := VerifyToken(c); err == nil {
		fmt.Println(cert)

		pagestr := c.Query("page")
		limitstr := c.Query("limit")
		appidstr := c.Query("appid")

		var tx *gorm.DB = PGDB
		if appidstr != "" {
			tx = PGDB.Where("appid = ?", appidstr)
		}

		var versions []Version
		if result := tx.Debug().Order("id").Find(&versions); result.Error == nil {

			items := versions
			total := int64(len(versions))

			if pagestr != "" {
				limitnum := convert.StoI(limitstr)
				startnum := (convert.StoI(pagestr) - 1) * limitnum
				endnum := convert.StoI(pagestr) * limitnum
				if endnum > total {
					endnum = total
				}
				items = versions[startnum:endnum]
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

func queryVersion(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	appidstr := c.Query("appid")
	versionstr := c.Query("version")

	// 只要验证信息通过就能查询
	if _, err := VerifyToken(c); err == nil {
		var version Version
		if result := PGDB.Debug().Where("appid = ? and version = ?", appidstr, versionstr).First(&version); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": version})
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
