package router

import (
	"bytes"
	"drugims/controller"
	"drugims/middleware"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 跨域访问处理
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问，也可设置特定域名
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Token, token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Origin, Content-Type, Token, token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// printRequestBodyMiddleware 中间件用于打印请求体
func printRequestBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取请求体
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("读取请求体失败:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 将请求体内容打印到控制台
		fmt.Println("请求体内容:", string(body))

		// 重新设置请求体，以便后续处理可以再次读取
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// 继续处理请求
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 设置跨域访问处理、debug中间件
	r.Use(corsMiddleware(), printRequestBodyMiddleware())

	// Home主页
	r.GET("/home", controller.Home)
	// 用户路由组
	userGroup := r.Group("/user")
	{
		// 用户注册
		userGroup.POST("/register", controller.UserRegister)
		// 用户登录
		userGroup.POST("/login", controller.UserLogin)
		// 更新用户信息
		userGroup.POST("/update", middleware.JWTMiddleWare(), controller.UserUpdate)
		// 获取用户列表
		userGroup.POST("/getUserList", middleware.JWTMiddleWare(), controller.GetUserList)
		// 用户注销(软删除)
		userGroup.POST("/delete", middleware.JWTMiddleWare(), controller.UserDelete)
	}

	// 药品路由组
	drugGroup := r.Group("/drug")
	{
		// 获取全部药品
		drugGroup.POST("/getDrugList", middleware.JWTMiddleWare(), controller.GetDrugList)
		// // 创建药品
		// drugGroup.POST("/create")
		// // 删除药品
		// drugGroup.POST("/delete")
		// // 购买药品
		// drugGroup.POST("/buy")
		// // 进货
		// drugGroup.POST("/provide")
	}
	return r
}
