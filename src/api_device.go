package main

import (
	"encoding/json"
	"fmt"
	"libs/convert"
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
	case "query-device":
		queryDevice(c)
	case "query-device-system":
		queryDeviceSystem(c)
	case "create-device":
		createDevice(c)
	case "update-device":
		updateDevice(c)
	case "delete-device":
		deleteDevice(c)
	case "push-update":
		pushUpdateToDevice(c)
	case "generate-jwt":
		generatejwt(c)
	case "verify-jwt":
		verifyjwt(c)
	case "login":
	case "logout":
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
	if account, err := VerifyTokenForAccount(c); err == nil {
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

func queryDevice(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	ouidstr := c.Query("ouid")
	// 只要验证信息通过就能查询设备信息？,不是
	if cert, err := VerifyToken(c); err == nil {
		if cert.CertType == "account" {
			// acc, _ := cert.CertContent.(Account)

		} else if cert.CertType == "device" {
			dev, _ := cert.CertContent.(Device)
			if dev.OUID != ouidstr {
				c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问"})
				return
			}
		} else {
			// 这个肯定不行，没有权限
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问"})
			return
		}

		var device Device
		if result := PGDB.Debug().Preload("System").Preload("DeviceStatus").Preload("License").Preload("OwnerCompany").Where("ouid = ?", ouidstr).First(&device); result.Error == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": device})
			return
		} else {
			c.JSON(200, gin.H{"errcode": 10201, "errmsg": "数据错误", "data": nil})
			return
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
		return
	}
}

func queryDeviceSystem(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}

	ouidstr := c.Query("ouid")
	if _, err := VerifyToken(c); err == nil {
		var device Device
		if result := PGDB.Debug().Preload("System").Where("ouid = ?", ouidstr).First(&device); result.Error == nil {

			applist := make([]map[string]interface{}, 0)
			appliststr := "" // 这个是加入到applist的appid字符串组合
			logger.Debugf("设备信息:%v", device)
			listStr := ""
			if device.System != nil {
				listStr = device.System.List
			}
		LOOP:
			appArray := strings.Split(listStr, ",")
			listStr = ""
			content := make(map[string]string)
			for i := 0; i < len(appArray); i++ {
				if !strings.Contains(appliststr, appArray[i]) {
					var app Application
					if result := PGDB.Where("appid = ?", appArray[i]).First(&app); result.Error == nil {
						appMap := make(map[string]interface{})
						appMap["app_id"] = app.Appid
						appMap["app_name"] = app.Name
						appMap["app_type"] = app.Type
						appMap["app_latest"] = app.Latest
						appMap["app_depend"] = app.Depend
						appMap["app_status"] = app.Status
						appMap["app_remark"] = app.Remark
						content[app.Appid] = app.Latest
						applist = append(applist, appMap)
						appliststr += appArray[i] + ","
						if app.Depend != "" {
							listStr += app.Depend + ","
						}
					}
				}
			}
			if listStr != "" {
				goto LOOP
			}
			if device.System != nil {
				// 只显示设备系统
				data := make(map[string]interface{})
				data["system_ouid"] = device.System.OUID
				data["system_name"] = device.System.Name
				data["system_content"] = content
				data["system_applist"] = applist
				data["system_main"] = device.System.Main
				data["system_status"] = device.System.Status
				data["system_remark"] = device.System.Remark
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": data})
			} else {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": ""})
			}
			return
		} else {
			c.JSON(200, gin.H{"errcode": 10201, "errmsg": "数据错误", "data": nil})
			return
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
		return
	}
}

// 创建账号
func createDevice(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	if account, err := VerifyTokenForAccount(c); err == nil {
		// fmt.Println(account)
		PGDB.First(&account)
		fmt.Println("数据库查询了账户数据", account)
		if account.TypeID == 6 || account.TypeID == 7 || account.Status == 0 {
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
			var license License
			if licenseCode != "" {
				// 判断注册码权限
				device.LicenseCode = &licenseCode
				if result := PGDB.Where("code = ?", licenseCode).First(&license); result.Error == nil {
					// 查询到了
					var systemouid = ""
					if device.SystemOUID != nil {
						systemouid = *device.SystemOUID
					}

					if !(license.ExpiresAt > time.Now().Unix() && license.UseCount < license.Count && strings.Contains(license.Permit, systemouid)) {
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
			device.OwnerID = account.CompanyID
			device.Manufacturer = account.CompanyID
			if result := PGDB.Debug().Create(&device); result.Error == nil {
				// 创建成功后，注册码-1
				if licenseCode != "" {
					license.UseCount++
					PGDB.Save(&license)
				}
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

// 生成JWT
func generatejwt(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	userParam := make(map[string]interface{})
	if err := c.BindJSON(&userParam); err == nil {
		ouidStr, _ := userParam["ouid"].(string)
		pinStr, _ := userParam["pin"].(string)
		var device Device
		if result := PGDB.Where("ouid = ?", ouidStr).First(&device); result.Error == nil {
			passwd := SHA256(device.PIN + "woshiyanzhi") // 将提交的密码加上盐值（woshiyanzhi）后哈希对比
			if pinStr == passwd {
				othData := make(map[string]interface{})
				othData["cert_type"] = "account"
				othData["ouid"] = device.OUID
				othData["name"] = device.Name
				othData["ip"] = c.ClientIP()
				othByte, _ := json.Marshal(othData)

				jwtheader := ouid.JWTHeader{Typ: "proof", Alg: "hs256"}
				aUnix := time.Now().Unix()
				expUnix := aUnix + 30*24*3600 // 设置30天有效期
				jwtproof := ouid.JWTProof{
					Jti: MD5(time.Now().UTC().Format("20060102150405")),
					Iss: Config.Ouid,
					Sub: device.OUID,
					Obj: device.OUID,
					Iat: convert.StoI(convert.UNIX2UTC(aUnix)),
					Nbf: convert.StoI(convert.UNIX2UTC(aUnix)),
					Exp: convert.StoI(convert.UNIX2UTC(expUnix)),
					Oth: string(othByte),
				}

				jwt := ouid.JWT{Header: jwtheader, Playload: jwtproof}
				jwthex, err := ouid.SignJWT(jwt, Config.PriKey)
				if err != nil {
					logger.Errorf("[%s]: %v", "签发JWT", err)
				}
				redata := make(map[string]interface{})
				redata["token"] = jwthex
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": redata})
			} else {
				// 验证不通过
				c.JSON(200, gin.H{"errcode": 10102, "errmsg": "密码错误"})
			}

		} else {
			// 请求数据错误
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求数据错误"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
	}
}

func verifyjwt(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

	userParam := make(map[string]interface{})
	if err := c.BindJSON(&userParam); err == nil {
		tokenStr, _ := userParam["token"].(string)
		// 读取jwt
		if jwtok, err := ouid.VerifyJWT(tokenStr, Config.PriKey); jwtok && err == nil {
			c.JSON(200, gin.H{"errcode": 0, "errmsg": "验证密钥成功", "data": 1})
		} else {
			c.JSON(200, gin.H{"errcode": 10107, "errmsg": "验证密钥失败"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
	}
}
