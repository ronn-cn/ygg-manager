package main

import (
	"encoding/json"
	"errors"
	"libs/ouid"
	"strings"

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
					acc.OUID, _ = accMap["ouid"].(string)
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
	case "system":
		handleSystem(c, ps)
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
