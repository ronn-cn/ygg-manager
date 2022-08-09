package main

import (
	"encoding/json"
	"fmt"
	"libs/convert"
	"libs/logger"
	"libs/ouid"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理账号接口
func handleAccount(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "login":
		accountLogin(c)
	case "logout":
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功"})
		return
	case "get-account-list":
		getAccountList(c)
	case "get-account-info":
		getAccountInfo(c)
	case "get-account-permission":
		getAccountPermissions(c)
	case "create-account":
		createAccount(c)
	case "update-account":
		updateAccount(c)
	case "delete-account":
		deleteAccount(c)
	default:
		c.Status(404)
		return
	}
}

func getAccountList(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyTokenForAccount(c); err == nil {
		fmt.Println(account)
		// TODO: 暂时设置为9，方便调用，正确改为2
		if account.TypeID > 2 { // 不是管理员，没有权限访问
			c.JSON(200, gin.H{"errcode": 10104, "errmsg": "没有权限访问"})
			return
		}
		// if account.Status == 0 { // 禁用状态
		// 	c.JSON(200, gin.H{"errcode": 10101, "errmsg": "账号禁用状态"})
		// 	return
		// }

		// 查询数据库操作
		var accounts []AccountApi
		if account.TypeID == 1 {
			PGDB.Preload("PermissionType").Order("id").Find(&accounts)
		} else {
			PGDB.Where("type != 1").Preload("PermissionType").Order("id").Find(&accounts)
		}
		fmt.Println(gin.H{"total": len(accounts), "items": accounts})
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": gin.H{"total": len(accounts), "items": accounts}})
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

func getAccountInfo(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405) //请求中的方法被禁用
		return
	}
	if account, err := VerifyTokenForAccount(c); err == nil {
		PGDB.Where("account = ?", account.Account).First(&account)
		if account.Status == 0 { // 禁用状态
			c.JSON(200, gin.H{"errcode": 10101, "errmsg": "账号禁用状态"})
			return
		}
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": account})
	} else {
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

func getAccountPermissions(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Status(405)
		return
	}
	if account, err := VerifyTokenForAccount(c); err == nil {
		redata := make(map[string]interface{})
		PGDB.Where("account = ?", account.Account).Preload("PermissionType").First(&account)
		redata["id"] = account.PermissionType.ID
		redata["name"] = account.PermissionType.Name
		newMenus := make([]Menu, 0)
		var menus []Menu
		PGDB.Order("id").Find(&menus)
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
		c.JSON(200, gin.H{"errcode": 10105, "errmsg": "请求密钥错误"})
	}
}

// 创建账号
func createAccount(c *gin.Context) {
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
		logger.Debugf("提交的密码：%v", account.Passwd)
		account.Passwd = SHA256(account.Passwd + "woshiyanzhi")
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
}

// 更新账号
func updateAccount(c *gin.Context) {
	if c.Request.Method != "PUT" {
		c.Status(405)
		return
	}
	var account Account
	if err := c.BindJSON(&account); err == nil {
		account.IP = c.ClientIP()
		if account.Passwd != "" { // 修改密码操作
			if result := PGDB.Debug().Model(&Account{OUID: account.OUID}).Update("passwd", SHA256(account.Passwd+"woshiyanzhi")); result.Error == nil {
				account.Passwd = ""
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": account})
			} else {
				c.JSON(200, gin.H{"errcode": 10204, "errmsg": "更新数据错误"})
			}
		} else {
			if result := PGDB.Debug().Model(&Account{OUID: account.OUID}).Select("Name", "Account", "TypeID", "Contact", "Detail", "Status", "CompanyID").Updates(&account); result.Error == nil {
				c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": account})
			} else {
				c.JSON(200, gin.H{"errcode": 10204, "errmsg": "更新数据错误"})
			}
		}
	} else {
		logger.Debugf("%v", err)
		c.JSON(200, gin.H{"errcode": 10103, "errmsg": "请求参数错误"})
	}
}

// 删除账号
func deleteAccount(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.Status(405)
		return
	}
	// 删除账号
	ouidstr := c.Query("ouid")
	logger.Debugf("请求的OUID：%v", ouidstr)
	if result := PGDB.Debug().Where("ouid = ?", ouidstr).Delete(&Account{}); result.Error == nil {
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": ouidstr})
	} else {
		fmt.Println("数据库返回的结构:", result)
		c.JSON(200, gin.H{"errcode": 10205, "errmsg": "删除数据错误"})
	}
}

// 账号登录
func accountLogin(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}
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
			passwd := passwdStr + "woshiyanzhi"   // 将提交的密码加上盐值（woshiyanzhi）后哈希对比
			if account.Passwd == SHA256(passwd) { // 判断密码正确
				othData := make(map[string]interface{})
				othData["cert_type"] = "account"
				othData["ouid"] = account.OUID
				othData["account"] = accountStr
				othData["name"] = account.Name
				othData["type"] = account.TypeID
				othData["ip"] = c.ClientIP()
				othByte, _ := json.Marshal(othData)

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

				// 记录日志
				var record Record
				record.Type = 0 // 0表示账号 1表示设备
				record.Level = "INFO"
				record.Action = "登录"
				record.OUID = account.Account
				record.Info = account.Account + "账号登录成功"
				PGDB.Debug().Create(&record)
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
