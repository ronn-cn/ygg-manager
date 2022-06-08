package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"libs/convert"
	"libs/logger"
	"libs/ouid"
	"net"
	"time"

	N "libs/network"
)

var parentInfo ParentInfo
var SelfNode *N.Node
var GatherNode *N.Node

// 父级信息
type ParentInfo struct {
	ouid       string
	address    string
	netname    string
	managerUrl string
	otaUrl     string
}

// 查询节点信息
func queryParentInfo(addr string) {
	if conn, err := net.Dial("tcp", addr); err == nil {
		obj := N.NewOMTPDataToJSON("Query", N.REPORT, 1, nil)
		sendBytes, _ := N.NewOMTPMessage("json", obj)
		if err := N.ConnSend(conn, sendBytes); err == nil {
			// 发送成功 等待接收
			if recvBytes, err := N.ConnRecv(conn); err == nil {
				if omtp, err := N.ParseOMTP(recvBytes); err == nil {
					if omtp.Format == "json" {
						omtpjson := omtp.Data.(N.OMTPDataToJSON)
						if omtpjson.Type == N.NOTICE && omtpjson.Action == "RegularHost" {
							if omtpjsonDataMap, ok := omtpjson.Data.(map[string]interface{}); ok {
								parentInfo.ouid, _ = omtpjsonDataMap["ouid"].(string)
								parentInfo.address, _ = omtpjsonDataMap["address"].(string)
								parentInfo.netname, _ = omtpjsonDataMap["netname"].(string)
								parentInfo.managerUrl, _ = omtpjsonDataMap["manager_url"].(string)
								parentInfo.otaUrl, _ = omtpjsonDataMap["ota_url"].(string)
								logger.Infof("成功查询到节点信息")
							}
						}
					}
				} else {
					logger.Errorf("[%s]: %v", "本机", err)
				}
			} else {
				logger.Errorf("[%s]: %v", "本机", err)
			}
		} else {
			logger.Errorf("[%s]: %v", "本机", err)
		}
	} else {
		logger.Errorf("[%s]: %v", "本机", err)
	}
}

// 交换密钥
func exchangeKeys(gather *N.Node) {
	if gather == nil {
		logger.Errorf("[%s]: %v", "交换密钥", "Gather节点为空")
		return
	}
	if gather.Session != nil {
		ecc := N.NewEllipticECDH(elliptic.P256())
		if privKey1, pubKey1, err := ecc.GenerateKey(rand.Reader); err == nil {
			obj := N.NewOMTPDataToJSON("ECDH", N.REPORT, 1, N.MAP{"pubkey": ecc.MarshalA(pubKey1)})
			if msgBytes, err := N.NewOMTPMessage("json", obj); err == nil {
				if err := SelfNode.NodeSend(gather, msgBytes); err != nil {
					// 发送自身公钥
					logger.Errorf("[%s]: %v", gather.OUID, err)
					return
				}
			} else {
				logger.Errorf("[%s]: %v", gather.OUID, err)
				return
			}
			// 接收新消息
			if recvdata, err := gather.NodeRecv(); err == nil {
				if omtp, err := N.ParseOMTP(recvdata); err == nil {
					if omtp.Format == "json" {
						omtpjson := omtp.Data.(N.OMTPDataToJSON)
						if omtpjson.Type == N.NOTICE {
							switch omtpjson.Action {
							case "ECDH":
								if _, ok := omtpjson.Data.(map[string]interface{}); ok {
									obj := omtpjson.Data.(map[string]interface{})
									pubkeyStr, _ := obj["pubkey"].(string)
									verifykey0, _ := obj["verifykey"].(string)
									if pubKey2, ok := ecc.UnmarshalA(pubkeyStr); ok {
										if sharekey, err := ecc.GenerateSharedSecret(privKey1, pubKey2); err == nil {
											secret := hex.EncodeToString(sharekey)
											if verifykeyBytes, err := N.Encrypt(sharekey, sharekey); err == nil {
												verifykey := hex.EncodeToString(verifykeyBytes)
												if verifykey == verifykey0 {
													// 获得协商后的共享密钥
													gather.Session.Secret = secret
												}
											} else {
												logger.Errorf("[%s]: %v", gather.OUID, err)
											}
										} else {
											logger.Errorf("[%s]: %v", gather.OUID, err)
										}
									}
								}
							}
						}
					}
				} else {
					logger.Errorf("[%s]: %v", gather.OUID, err)
				}
			} else {
				logger.Errorf("[%s]: %v", gather.OUID, err)
			}
		} else {
			logger.Errorf("[%s]: %v", gather.OUID, err)
		}
	} else {
		logger.Errorf("[%s]: %v", "密钥交换", "gather.Session为空")
	}
}

// 申请加入
func applyJoin(gather *N.Node) {
	if gather.Session != nil && gather.Session.Secret != "" {
		hexjwt, err := GenerateJWT(Config.Ouid, Config.Pin, 5*60)
		if err != nil {
			logger.Errorf("[%s]: %v", "生成JWT令牌失败，详情", err)
		}
		addr := gather.Session.Conn.LocalAddr().String()
		obj := N.NewOMTPDataToJSON("Join", N.REPORT, 1, N.MAP{"ouid": Config.Ouid, "hexjwt": hexjwt, "address": addr})
		if msgBytes, err := N.NewOMTPMessage("json", obj); err == nil {
			if err = SelfNode.NodeSend(gather, msgBytes); err != nil {
				logger.Errorf("[%s]: %v", gather.OUID, err)
				return
			}
		} else {
			logger.Errorf("[%s]: %v", gather.OUID, err)
			return
		}
		gather.Relation = N.GATHER
		gather.Status = N.ONLINE
	}
}

// GenerateJWT 生成JWT
func GenerateJWT(ouidstr string, pinstr string, limit int64) (string, error) {
	var jwtheader ouid.JWTHeader
	jwtheader.Typ = "proof"
	jwtheader.Alg = "hs256"
	var jwtproof ouid.JWTProof
	jwtproof.Jti = MD5(time.Now().UTC().Format("20060102150405"))
	jwtproof.Iss = ouidstr
	jwtproof.Sub = ouidstr
	jwtproof.Obj = ouidstr
	aUnix := time.Now().Unix()
	jwtproof.Iat = convert.StoI(convert.UNIX2UTC(aUnix))
	jwtproof.Nbf = convert.StoI(convert.UNIX2UTC(aUnix))
	expUnix := aUnix + limit
	jwtproof.Exp = convert.StoI(convert.UNIX2UTC(expUnix))

	var jwt ouid.JWT
	jwt.Header = jwtheader
	jwt.Playload = jwtproof

	return ouid.SignJWT(jwt, MD5(pinstr))
}

// ConnGather 加入到Gather节点网络中
func connGather(ouidstr string, addr string) *N.Node {
	gather := N.NewNode(ouidstr)
	if conn, err := net.Dial("tcp", addr); err != nil {
		logger.Errorf("[%s]: %v", "连接父级地址失败，详情", err)
		return nil
	} else {
		gather.Session = N.NewNodeSession(conn)
		exchangeKeys(gather) // 交换加密密钥
		applyJoin(gather)    // 申请加入父级网络
	}
	logger.Infof("成功加入到Yggdrasil节点")

	// 处理父级节点
	go handleGatherNode(gather)
	return gather
}

// 处理父级节点
func handleGatherNode(gather *N.Node) {
	if gather.Relation != N.GATHER {
		logger.Warnf("[%s]: %v", gather.OUID, "不是GATHER节点")
		gather = nil
		return
	}
	for {
		if _, err := gather.NodeRecv(); err != nil {
			logger.Errorf("[%s]: %v", gather.OUID, err)
			gather = nil
			return
		}
	}
}

// 连接Yggdrasil
func connYggdrasil() {
	if SelfNode == nil {
		SelfNode = N.NewNode(Config.Ouid)
	}
	for GatherNode == nil {
		if Config.Yggdrasil != "" {
			queryParentInfo(Config.Yggdrasil)
			GatherNode = connGather(parentInfo.ouid, parentInfo.address)
		}
		time.Sleep(5 * time.Second)
	}
}
