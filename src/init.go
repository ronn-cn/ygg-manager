package main

import (
	"fmt"

	"libs/logger"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Config Configure
var PGDB *gorm.DB // 数据库

// // go:embed view/*
// var viewFS embed.FS

// 初始化配置
func initConfig() bool {
	var tomlPath = "./conf/manager.toml"
	if tree, err := toml.LoadFile(tomlPath); err == nil {
		if err := tree.Unmarshal(&Config); err == nil {
			logger.Infof("初始化配置成功")
			return true
		} else {
			logger.Errorf("生成结构体错误：%v", err)
		}
	} else {
		logger.Errorf("加载配置文件错误：%v", err)
	}
	return false
}

// 初始化日志
func initLogger() bool {
	logger.Config("./log/manager.log", logger.DebugLevel)
	logger.SetSaveMode(logger.ByDay) // 设置保存类型
	logger.SetSaveDays(30)           // 设置保存天数
	logger.SetDebug(true)            // 设置调试

	logger.Infof("初始化日志完成")
	return true
}

// 初始化数据库
func initDB() bool {
	host := "172.16.7.100"
	port := "5432"
	user := "postgres"
	passwd := "postgres"
	dbname := "manager"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, passwd, dbname, port)
	var pgErr error
	PGDB, pgErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if pgErr != nil {
		logger.Errorf("连接数据库错误：%v", pgErr)
		return false
	}
	logger.Infof("初始化数据库完成")
	fmt.Println(PGDB)
	return true
}

func initServer() {
	port := "8800"
	if Config.Port != "" {
		port = Config.Port
	}
	router := gin.Default()
	initRouter(router) //初始化路由
	logger.Infof("初始化Web服务完成")
	if err := router.Run(":" + port); err != nil {
		logger.Errorf("启动Web服务错误：%v", err)
	}
}
