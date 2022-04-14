package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 处理设备接口
func handleDevice(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-device-list":
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		if account, err := VerifyToken(c); err == nil {
			fmt.Println(account)

			// 查询数据库操作
			var devices []Device
			switch account.TypeID {
			case 3: // 判断是开发
				fmt.Println("判断是开发,只能查看开发设备")
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("status = ?", 13).Find(&devices)
			case 4: // 判断是运维
				fmt.Println("判断是运维,只能查看除开发之外的设备")
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("status <> ?", 13).Find(&devices)
			case 5: // 判断是生产
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("manufacturer = ?", account.Company).Find(&devices)
			case 6: // 判断是销售
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("distributor = ?", account.Company).Find(&devices)
			case 7: // 判断是客户
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("customer = ?", account.Company).Find(&devices)
			default:
				PGDB.Preload("Assembly").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Find(&devices)
			}
			fmt.Println(gin.H{"total": len(devices), "items": devices})
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(devices), "items": devices}})
		} else {
			c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
		}

	default:
		c.Status(404)
		return
	}
}
