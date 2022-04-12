package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"libs/logger"
	"libs/ouid"
	"strings"
	"time"

	"libs/convert"

	"github.com/gin-gonic/gin"
)

// 验证签名
func VerifyToken(c *gin.Context) (Account, error) {
	var acc Account
	jwtoken := c.GetHeader("Authorization")
	if jwtoken == "" {
		return acc, errors.New("没有验证信息")
	} else {
		authArr := strings.Split(jwtoken, " ")
		if len(authArr) != 2 || authArr[0] != "hexjwt" {
			return acc, errors.New("验证信息格式不正确")
		}
		// 读取jwt
		if jwtok, err := ouid.VerifyJWT(authArr[1], Config.PriKey); jwtok && err == nil {
			if jwtObj, err := ouid.ParseJWT(authArr[1]); err == nil {
				playload := jwtObj.Playload.(ouid.JWTProof)
				accMap := make(map[string]interface{})
				if err := json.Unmarshal([]byte(playload.Oth), &accMap); err == nil {
					acc.Account, _ = accMap["account"].(string)
					acc.Name, _ = accMap["name"].(string)
					typeid, _ := accMap["type"].(float64)
					acc.TypeID = int(typeid)
					acc.IP, _ = accMap["ip"].(string)
					// logger.Debugf("%v", acc)
					return acc, nil
				} else {
					return acc, err
				}
			} else {
				return acc, err
			}
		} else {
			return acc, err
		}
	}
}

// 处理管理API
func handleManagerApi(c *gin.Context) {
	path := c.Param("path")
	ps := strings.Split(path[1:], "/")
	switch ps[0] {
	case "dashboard":
		handleDashboard(c, ps)
	case "account":
		handleAccount(c, ps)
	case "device":
		handleDevice(c, ps)
	case "assembly":
		handleAssembly(c, ps)
	case "application":
		handleApplication(c, ps)
	case "license":
		handleLicense(c, ps)
	case "company":
		handleCompany(c, ps)
	case "record":
		handleRecord(c, ps)
	case "setting":
		handleSetting(c, ps)
	default:
		c.Status(404)
	}
}

// 处理账号接口
func handleAccount(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-account-list":
		getAccountList(c)
	case "post-account":
		if c.Request.Method != "POST" {
			c.Status(405)
			return
		}
		// 创建账号
		var account Account
		if err := c.BindJSON(&account); err == nil {
			account.OUID = ouid.GenerateOUID()
			account.Status = 1
			account.IP = c.ClientIP()
			if result := PGDB.Debug().Create(&account); result.Error == nil {
				fmt.Println("1:", result)
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": account.OUID})
			} else {
				fmt.Println("2:", result)
				c.JSON(200, gin.H{"errcode": 10203, "errmsg": "创建数据错误"})
			}
		} else {
			logger.Debugf("%v", err)
			c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
		}
	case "get-account-info":
		getAccountInfo(c)
	case "get-account-permission":
		getAccountPermissions(c)
	case "login":
		login(c)
	case "logout":
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功"})
		return
	default:
		c.Status(404)
		return
	}
}

// 处理仪表板请求
func handleDashboard(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-data":
		reData := make(map[string]interface{})
		reData["device_count"] = 10
		reData["assembly_count"] = 11
		reData["application_count"] = 12
		reData["customer_count"] = 6
		reData["dealer_count"] = 5
		reData["facturer_count"] = 5
		reData["today_device_add"] = 5
		reData["today_assembly_add"] = 5
		reData["today_application_add"] = 5
		reData["today_device_produced"] = 5
		reData["today_device_sold"] = 5
		reData["today_device_installed"] = 5
		reData["account_count"] = 5
		reData["record_count"] = 5
		reData["today_record_add"] = 5
		reData["version_count"] = 5
		reData["today_version_add"] = 5
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": reData})
	default:
		c.Status(404)
	}
}

func getAccountList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyToken(c); err == nil {
		fmt.Println(account)
		// TODO: 暂时设置为9，方便调用，正确改为2
		if account.TypeID > 2 { // 不是管理员，没有权限访问
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问", "data": nil})
			return
		}
		// if account.Status == 0 { // 禁用状态
		// 	c.JSON(200, gin.H{"errcode": 10101, "errmsg": "账号禁用状态", "data": nil})
		// 	return
		// }

		// 查询数据库操作
		var accounts []Account
		if account.TypeID == 1 {
			PGDB.Preload("PermissionType").Find(&accounts)
		} else {
			PGDB.Where("type != 1").Preload("PermissionType").Find(&accounts)
		}
		fmt.Println(gin.H{"total": len(accounts), "items": accounts})
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(accounts), "items": accounts}})
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

func getAccountInfo(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405) //请求中的方法被禁用
		return
	}
	if account, err := VerifyToken(c); err == nil {
		PGDB.Where("account = ?", account.Account).First(&account)
		if account.Status == 0 { // 禁用状态
			c.JSON(200, gin.H{"errcode": 10101, "errmsg": "账号禁用状态", "data": nil})
			return
		}
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": account})
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

func getAccountPermissions(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyToken(c); err == nil {
		redata := make(map[string]interface{})
		PGDB.Where("account = ?", account.Account).Preload("PermissionType").First(&account)
		redata["id"] = account.PermissionType.ID
		redata["name"] = account.PermissionType.Name
		newMenus := make([]Menu, 0)
		var menus []Menu
		PGDB.Find(&menus)
		if account.PermissionType.AllowMenu == "0" {
			// 全部
			redata["menu"] = menus
		} else {
			menuIdArr := strings.Split(account.PermissionType.AllowMenu, ",")
			for _, value := range menuIdArr {
				for _, menuItem := range menus {
					if convert.ToS(menuItem.ID) == value {
						newMenus = append(newMenus, menuItem)
					}
				}
			}
			redata["menu"] = newMenus
		}
		fmt.Println(redata)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": redata})
	} else {
		fmt.Println(err)
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误", "data": nil})
	}
}

func login(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}
	// TODO:还需要做登录日志记录
	userParam := make(map[string]interface{})
	if err := c.BindJSON(&userParam); err == nil {
		accountStr, _ := userParam["account"].(string)
		passwdStr, _ := userParam["password"].(string)
		var account Account
		PGDB.Preload("PermissionType").Where("account = ?", accountStr).First(&account)
		if account.Account != "" {
			// fmt.Println(account.Status)
			if account.Status != 1 { // 账号处于被禁用状态
				c.JSON(200, gin.H{"errcode": 10101, "errmsg": "账号禁用"})
				return
			}
			passwd := passwdStr + "woshishen"     // 将提交的密码加上盐值（woshishen）后哈希对比
			if account.Passwd == SHA256(passwd) { // 判断密码正确
				othData := make(map[string]interface{})
				othData["account"] = accountStr
				othData["name"] = account.Name
				othData["type"] = account.TypeID
				othData["ip"] = c.ClientIP()
				othByte, _ := json.Marshal(othData)
				fmt.Println(string(othByte))

				jwtheader := ouid.JWTHeader{Typ: "proof", Alg: "hs256"}
				aUnix := time.Now().Unix()
				expUnix := aUnix + 30*24*3600 // 设置30天有效期
				jwtproof := ouid.JWTProof{
					Jti: MD5(time.Now().UTC().Format("20060102150405")),
					Iss: Config.Ouid,
					Sub: account.OUID,
					Obj: account.Account,
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
				c.JSON(200, gin.H{"errcode": 10102, "errmsg": "密码错误"})
			}
		} else {
			c.JSON(200, gin.H{"errcode": 10100, "errmsg": "账号不存在"})
		}
	} else {
		c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
	}
}

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

func handleAssembly(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-assembly-list": // 查询系统列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var assemblies []Assembly
		PGDB.Find(&assemblies)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(assemblies), "items": assemblies}})
	default:
		c.Status(404)
		return
	}
}

func handleApplication(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-application-list": // 查询应用列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var application []Application
		PGDB.Find(&application)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(application), "items": application}})
	default:
		c.Status(404)
		return
	}
}

func handleLicense(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-license-list": // 查询应用列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var license []License
		PGDB.Find(&license)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(license), "items": license}})
	default:
		c.Status(404)
		return
	}
}

func handleCompany(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-company-list": // 查询应用列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var companies []Company // 公司复数
		PGDB.Find(&companies)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(companies), "items": companies}})
	default:
		c.Status(404)
		return
	}
}

func handleRecord(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-record-list": // 查询应用列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var records []Record // 日志复数
		PGDB.Find(&records)
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(records), "items": records}})
	default:
		c.Status(404)
		return
	}
}

func handleSetting(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-setting": // 查询应用列表
		if c.Request.Method != "GET" {
			c.Status(405)
			return
		}
		var settings []Setting // 设置复数
		PGDB.Find(&settings)
		redata := make(map[string]interface{})
		for _, item := range settings {
			redata[item.Key] = item.Value
		}
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": redata})
	default:
		c.Status(404)
		return
	}
}
