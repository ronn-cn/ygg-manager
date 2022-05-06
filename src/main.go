package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	defer fmt.Println("Exit!")
	if ok := initLogger(); !ok {
		return
	}
	if ok := initConfig(); !ok {
		return
	}
	if ok := initDB(); !ok {
		return
	}
	// 将配置写入数据库
	b, _ := json.Marshal(&Config)
	var conf map[string]interface{}
	_ = json.Unmarshal(b, &conf)
	for k, v := range conf {
		// fmt.Printf("key:%v value:%v\n", k, v)
		var sItem Setting
		if result := PGDB.Where("key = ?", k).Find(&sItem); result.Error == nil {
			// 查询到了
			 
		} else {
			// 没有查询到
		}
	}
	initServer()
}
