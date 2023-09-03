package main

import (
	"chat-server/router"
	"chat-server/utils"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {

	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	utils.InitOss()
	fmt.Println("oss connect...", oss.Version)
	r := router.Router()
	r.Run(":8090")
}
