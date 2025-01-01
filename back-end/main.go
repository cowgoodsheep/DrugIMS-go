package main

import (
	"drugims/config"
	"drugims/dao"
	"drugims/model"
	"drugims/router"
	"fmt"
)

func main() {
	//数据库初始化
	dao.InitMySQL()
	//数据库迁移
	dao.DB.AutoMigrate(&model.UserInfo{})
	//关闭数据库
	defer dao.DB.Close()
	//开启路由
	r := router.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		fmt.Printf("server start failed, error:%v\n", err)
	}
}
