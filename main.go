package main

import (
	"chat-server/router"
	"chat-server/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8090")
}
