package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"project2/docs"
	"project2/service"
)

//地址映射

func Router() *gin.Engine {
	//创建路由
	r := gin.Default()
	docs.SwaggerInfo.Title = ""
	//swag测试界面
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("templates/**/*")

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)

	//用户模块
	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/register", service.CreateUser)
	r.GET("/logout", service.DeleteUser)
	r.POST("/updateInfo", service.UpdateUser)
	r.POST("findUserByNameAndPwd", service.FindUserByNameAndPwd)
	//发送消息
	r.GET("/sendMsg", service.SendMsg)
	r.GET("/sendUserMsg", service.SendUserMsg)
	return r
}
