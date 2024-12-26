package routers

import (
	"drugims/controller"
	"drugims/middleware"

	"github.com/gin-gonic/gin"
)

// 跨域访问处理
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问，也可设置特定域名
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 设置跨域访问处理中间件
	r.Use(corsMiddleware())

	//Home主页
	r.GET("/home", controller.Home)
	// 用户路由组
	userGroup := r.Group("/user")
	{
		//用户注册
		userGroup.POST("/register", middleware.SHAMiddleWare(), controller.UserRegister)
		//用户登录
		userGroup.POST("/login", middleware.SHAMiddleWare(), controller.UserLogin)
	}
	return r
}
