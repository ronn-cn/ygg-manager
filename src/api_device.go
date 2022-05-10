package main

import (
	"fmt"
	"libs/logger"
	"libs/ouid"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 处理设备接口
func handleDevice(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-device-list":
		getDeviceList(c)
	case "create-device":
		createDevice(c)
	case "update-device":
		updateDevice(c)
	case "delete-device":
		deleteDevice(c)
	case "push-update":
		pushUpdateToDevice(c)
	default:
		c.Status(404)
		return
	}
}

func getDeviceList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)

		namestr := c.Query("name")
		ouidstr := c.Query("ouid")
		statusstr := c.Query("status_id")
		systemstr := c.Query("system_ouid")

		var tx *gorm.DB = PGDB
		if namestr != "" {
			tx = PGDB.Where("name LIKE ?", "%"+namestr+"%")
		}
		if ouidstr != "" {
			tx = PGDB.Where("ouid LIKE ?", "%"+ouidstr+"%")
		}
		if statusstr != "" {
			tx = PGDB.Where("status = ?", statusstr)
		}
		if systemstr != "" {
			tx = PGDB.Where("system = ?", systemstr)
		}

		// 查询数据库操作
		var devices []Device
		switch account.TypeID {
		case 3: // 判断是开发
			fmt.Println("判断是开发,只能查看开发设备")
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("status = ?", 13).Find(&devices)
		case 4: // 判断是运维
			fmt.Println("判断是运维,只能查看除开发之外的设备")
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("status <> ?", 13).Find(&devices)
		case 5: // 判断是生产
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("manufacturer = ?", account.Company).Find(&devices)
		case 6: // 判断是销售
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("distributor = ?", account.Company).Find(&devices)
		case 7: // 判断是客户
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("customer = ?", account.Company).Find(&devices)
		default:
			tx.Order("id").Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Find(&devices)
		}
		fmt.Println(gin.H{"total": len(devices), "items": devices})
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(devices), "items": devices}})
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 创建账号
func createDevice(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		if account.TypeID == 6 || account.TypeID == 7 {
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问"})
			return
		}

		// 创建账号
		var device Device
		if err := c.BindJSON(&device); err == nil {
			if account.TypeID == 3 { //开发类型
				var id int = 13
				device.StatusID = &id
			}
			licenseCode := c.Query("license")
			if licenseCode != "" {
				// 判断注册码权限
				device.LicenseCode = &licenseCode
				var license License
				if result := PGDB.Where("code = ?", licenseCode).First(&license); result.Error == nil {
					// 查询到了
					if !(license.ExpiresAt > time.Now().Unix() && license.UseCount < license.Count && strings.Contains(license.Permit, *device.SystemOUID)) {
						c.JSON(200, gin.H{"errcode": 10104, "errmsg": "授权码无效，没有权限访问"})
						return
					}
				}
			}
			if device.OUID == "" {
				device.OUID = ouid.GenerateOUID()
			}
			if device.Name == "" {
				device.Name = "未命名的新设备"
			}
			if device.PIN == "" {
				rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
				device.PIN = fmt.Sprintf("%06v", rnd.Int31n(1000000))
			}
			if result := PGDB.Debug().Create(&device); result.Error == nil {
				// 创建成功后，注册码-1
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": device.OUID})
			} else {
				c.JSON(200, gin.H{"errcode": 10203, "errmsg": "创建数据错误"})
			}
		} else {
			logger.Debugf("%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 更新设备
func updateDevice(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		var device Device
		if err := c.BindJSON(&device); err == nil {
			if result := PGDB.Debug().Model(&Device{OUID: device.OUID}).Select("name", "pin", "system", "status", "remark").Updates(&device); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": device})
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

// 删除设备
func deleteDevice(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}

	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		// 删除账号
		ouidstr := c.Query("ouid")
		logger.Debugf("请求的OUID：%v", ouidstr)
		if result := PGDB.Debug().Where("ouid = ?", ouidstr).Delete(&Device{}); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": ouidstr})
		} else {
			fmt.Println("数据库返回的结构:", result)
			c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

// 推送更新
func pushUpdateToDevice(c *gin.Context) {
	// 调用YGG接口
}
