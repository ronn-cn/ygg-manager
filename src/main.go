package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
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
	configToDB()
	initServer()
}

// 配置文件写入数据库
func configToDB() {
	// 将配置写入数据库
	b, _ := json.Marshal(&Config)
	var conf map[string]string
	_ = json.Unmarshal(b, &conf)
	for k, v := range conf {
		var sItem Setting
		if result := PGDB.Where("key = ?", k).First(&sItem); result.Error == nil {
			if sItem.Value != v {
				sItem.Value = v
				PGDB.Save(&sItem)
			}
		} else {
			// 没有查询到
			sItem.Key = k
			sItem.Value = v
			PGDB.Create(&sItem)
		}
	}
}

func configToFile() {
	cfgBytes, _ := toml.Marshal(Config)
	filepath := "./conf/manager.toml"
	ioutil.WriteFile(filepath, cfgBytes, 0666)
}
