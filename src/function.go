package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// MD5 算法
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 算法
func SHA256(str string) string {
	bytes2 := sha256.Sum256([]byte(str))       //计算哈希值，返回一个长度为32的数组
	hashcode2 := hex.EncodeToString(bytes2[:]) //将数组转换成切片，转换成16进制，返回字符串
	return hashcode2
}
