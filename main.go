package main

import (
	"project2/router"
	"project2/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run() //修改端口号r.Run(":8081")，默认为8080
}
