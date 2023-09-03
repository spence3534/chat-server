package router

import (
	"chat-server/docs"
	"chat-server/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 登录
	r.POST("/login", service.Login)
	// 注册
	r.POST("/createUser", service.CreateUser)
	r.GET("/getUserList", service.GetUserList)

	// 首页相关
	home := r.Group("home")
	{
		home.GET("/index", service.GetIndex)
	}

	// 用户相关
	user := r.Group("user") // .Use(middleware.JwtAuth())
	{
		user.GET("/deleteUser", service.DeleteUser)
		user.POST("/updateUserInfo", service.UpdateUserInfo)
		user.GET("/getUserInfo", service.GetUserInfo)
		user.GET("/sendMessage", service.SendMessage)
		user.GET("/sendUserMsg", service.SendUserMsg)
		user.GET("/getFriendList", service.GetFriendList)
		user.GET("/dddFriend", service.AddFriend)
		user.GET("/searchFriend", service.SearchFriend)
	}

	// 上传文件
	file := r.Group("file")
	{
		file.POST("/uploadImage", service.UploadImage)
	}
	return r
}
