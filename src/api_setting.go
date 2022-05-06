package main

import "github.com/gin-gonic/gin"

func handleSetting(c *gin.Context, ps []string) {
	if len(ps) < 1 {
		c.Status(404)
		return
	}
	switch ps[1] {
	case "get-setting": // 查询应用列表
		getSetting(c)
	case "save-setting": //保存设置
	default:
		c.Status(404)
		return
	}
}

func getSetting(c *gin.Context) {
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
}

func saveSetting(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Status(405)
		return
	}

}
