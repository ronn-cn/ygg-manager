package main

import (
	"encoding/json"
	"errors"
	"libs/ouid"
	"strings"
	"time"

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
		// 查询设备总数
		var devices []Device
		var deviceCount int64 = 0
		var deviceTrend = []int{0, 0, 0, 0, 0, 0, 0}
		if result := PGDB.Debug().Find(&devices); result.Error == nil {
			deviceCount = result.RowsAffected
		}

		// 查询系统总数
		var systems []System
		var systemCount int64 = 0
		// var systemTrend = []int{0, 0, 0, 0, 0, 0, 0}
		if result := PGDB.Debug().Find(&systems); result.Error == nil {
			systemCount = result.RowsAffected
		}

		// 查询应用总数
		var apps []Application
		var appCount int64 = 0
		if result := PGDB.Debug().Find(&apps); result.Error == nil {
			appCount = result.RowsAffected
		}

		// 查询版本总数
		var versions []Version
		var versionCount int64 = 0
		if result := PGDB.Debug().Find(&versions); result.Error == nil {
			versionCount = result.RowsAffected
		}

		// 查询账户总数
		var accounts []Account
		var accountCount int64 = 0
		if result := PGDB.Debug().Find(&accounts); result.Error == nil {
			accountCount = result.RowsAffected
		}

		// 查询密钥总数
		var licenses []License
		var licenseCount int64 = 0
		if result := PGDB.Debug().Find(&licenses); result.Error == nil {
			licenseCount = result.RowsAffected
		}

		nTime := time.Now()
		var oldTimeUnix int64 = nTime.Unix()
		for i := 0; i < 7; i++ {
			yTime := nTime.AddDate(0, 0, -i)
			dateStr := yTime.Format("2006-01-02")
			dayTime, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
			dayTimeUnix := dayTime.Unix()

			for _, dev := range devices {
				if dev.CreatedAt >= dayTimeUnix && dev.CreatedAt < oldTimeUnix {
					deviceTrend[i]++
				}
			}
			oldTimeUnix = dayTimeUnix
		}

		reData := make(map[string]interface{})
		reData["device_count"] = deviceCount
		reData["system_count"] = systemCount
		reData["application_count"] = appCount
		reData["version_count"] = versionCount
		reData["customer_count"] = 6 // 客户总数
		reData["dealer_count"] = 5   // 经销商总数
		reData["facturer_count"] = 5 //生产商总数
		reData["account_count"] = accountCount
		reData["record_count"] = 5
		reData["license_count"] = licenseCount
		reData["today_device_add"] = 5
		reData["today_system_add"] = 5
		reData["today_application_add"] = 5
		reData["today_device_produced"] = 5
		reData["today_device_sold"] = 5
		reData["today_device_installed"] = 5
		reData["today_record_add"] = 5
		reData["today_version_add"] = 5
		reData["device_trend"] = deviceTrend
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
