package main

import (
	"encoding/json"
	"errors"
	"libs/ouid"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 证明
type Certification struct {
	OUID        string
	Name        string
	IP          string
	CertType    string
	CertContent interface{}
}

// 验证签名
func VerifyToken(c *gin.Context) (Certification, error) {
	var cert Certification
	jwtoken := c.GetHeader("Authorization")
	if jwtoken == "" {
		return cert, errors.New("没有验证信息")
	} else {
		authArr := strings.Split(jwtoken, " ")
		if len(authArr) != 2 || authArr[0] != "hexjwt" {
			return cert, errors.New("验证信息格式不正确")
		}
		// 读取jwt
		if jwtok, err := ouid.VerifyJWT(authArr[1], Config.PriKey); jwtok && err == nil {
			if jwtObj, err := ouid.ParseJWT(authArr[1]); err == nil {
				playload := jwtObj.Playload.(ouid.JWTProof)
				certMap := make(map[string]interface{})
				if err := json.Unmarshal([]byte(playload.Oth), &certMap); err == nil {
					certType, _ := certMap["cert_type"].(string)
					cert.CertType = certType
					cert.OUID, _ = certMap["ouid"].(string)
					cert.Name, _ = certMap["name"].(string)
					cert.IP, _ = certMap["ip"].(string)

					switch certType {
					case "account":
						var acc Account
						acc.OUID, _ = certMap["ouid"].(string)
						acc.Account, _ = certMap["account"].(string)
						acc.Name, _ = certMap["name"].(string)
						typeid, _ := certMap["type"].(float64)
						acc.TypeID = int(typeid)
						acc.IP, _ = certMap["ip"].(string)
						cert.CertContent = acc
					case "device":
					}

					return cert, nil
				} else {
					return cert, err
				}
			} else {
				return cert, err
			}
		} else {
			return cert, err
		}
	}
}

// 验证账户的令牌
func VerifyTokenForAccount(c *gin.Context) (Account, error) {
	var acc Account
	if cert, err := VerifyToken(c); err == nil {
		if cert.CertType == "account" {
			acc = cert.CertContent.(Account)
		} else {
			err = errors.New("不是账号类型的验证信息")
		}
		return acc, err
	} else {
		return acc, err
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
		var systemTrend = []int{0, 0, 0, 0, 0, 0, 0}
		if result := PGDB.Debug().Find(&systems); result.Error == nil {
			systemCount = result.RowsAffected
		}

		// 查询应用总数
		var apps []Application
		var appCount int64 = 0
		var appTrend = []int{0, 0, 0, 0, 0, 0, 0}
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

			for _, sys := range systems {
				if sys.CreatedAt >= dayTimeUnix && sys.CreatedAt < oldTimeUnix {
					systemTrend[i]++
				}
			}

			for _, app := range apps {
				if app.CreatedAt >= dayTimeUnix && app.CreatedAt < oldTimeUnix {
					appTrend[i]++
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
		reData["system_trend"] = systemTrend
		reData["app_trend"] = appTrend
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
