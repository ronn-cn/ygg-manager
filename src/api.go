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
		} // 读取jwt
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
						acc.OUID = cert.OUID
						acc.Account, _ = certMap["account"].(string)
						acc.Name = cert.Name
						typeid, _ := certMap["type"].(float64)
						acc.TypeID = int(typeid)
						acc.IP = cert.IP
						cert.CertContent = acc
					case "device":
						// TODO: 编写设备的验证,现在写了，但是还是缺少。
						var dev Device
						dev.OUID = cert.OUID
						dev.Name = cert.Name
						cert.CertContent = dev
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

// TODO:验证设备的令牌
func VerifyTokenForDevice(c *gin.Context) (Device, error) {
	var dev Device
	if cert, err := VerifyToken(c); err == nil {
		if cert.CertType == "device" {
			dev = cert.CertContent.(Device)
		} else {
			err = errors.New("不是设备类型的验证信息")
		}
		return dev, err
	} else {
		return dev, err
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
	case "version":
		handleVersion(c, ps)
	case "license":
		handleLicense(c, ps)
	case "company":
		handleCompany(c, ps)
	case "record":
		handleRecord(c, ps)
	case "setting":
		handleSetting(c, ps)
	case "bbs":
		handleBBS(c, ps)
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
		reData["customer_count"] = 0 // 客户总数
		reData["dealer_count"] = 0   // 经销商总数
		reData["facturer_count"] = 0 //生产商总数
		reData["account_count"] = accountCount
		reData["record_count"] = 0
		reData["license_count"] = licenseCount
		reData["today_device_add"] = 0
		reData["today_system_add"] = 0
		reData["today_application_add"] = 0
		reData["today_device_produced"] = 0
		reData["today_device_sold"] = 0
		reData["today_device_installed"] = 0
		reData["today_record_add"] = 0
		reData["today_version_add"] = 0
		reData["device_trend"] = deviceTrend
		reData["system_trend"] = systemTrend
		reData["app_trend"] = appTrend
		c.JSON(200, gin.H{"errcode": 0, "errmsg": "请求成功", "data": reData})
	default:
		c.Status(404)
	}
}

func CreatJwt(toouid string, registerCode string) (string, error) {

	jwt := ouid.JWT{
		Header: ouid.JWTHeader{
			Typ: "proof",
			Alg: "hs256",
		},
		Playload: ouid.JWTProof{
			Jti: MD5(time.Now().Local().String()),
			Iss: Config.Ouid,
			Sub: toouid,
			Obj: toouid,
			Iat: 0,
			Nbf: 0,
			Oth: registerCode,
		},
	}

	return ouid.SignJWT(jwt, Config.PriKey)

}
func handleBBS(c *gin.Context, ps []string) {
	// 记录，这是一个老接口，从原先manager1版本中移植过来的，用于bss的验证
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "jwt":
		if c.Request.Method == "GET" {
			if c.Query("enterInfo") == "" {
				c.JSON(200, gin.H{"code": 0, "msg": "缺少必要参数"})
				return
			}
			// 为bbs生成jwt
			// 签发jwt证书
			if jwtRsp, err := CreatJwt(c.Query("enterInfo"), ""); err == nil {
				c.JSON(200, gin.H{"code": 1, "msg": "获取成功", "data": jwtRsp})
				return
			}
			c.JSON(200, gin.H{"code": 0, "msg": "获取失败"})
		}

		if c.Request.Method == "POST" {
			token := c.PostForm("jwt")
			if token == "" {
				c.JSON(200, gin.H{"code": 0, "msg": "必要参数为空"})
				return
			}

			if _, err := ouid.ParseJWT(token); err == nil {

				c.JSON(200, gin.H{"code": 1, "msg": "解析成功"})
				return
			}
			c.JSON(200, gin.H{"code": 0, "msg": "解析失败"})
			return
		}
	}
}
