package controller

import (
	"drugims/logic"
	"drugims/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(c *gin.Context) {
	// 获取用户注册信息
	var u model.UserInfo
	u.UserName = c.PostForm("user_name")
	u.Telephone = c.PostForm("telephone")
	passwordTemp, ok := c.Get("password")
	u.Password = passwordTemp.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "密码解析错误",
		})
		return
	}

	// 上传用户注册信息至用户登录服务，进行用户注册
	registerMsg, err := logic.RegisterUser(u.Telephone, u.UserName, u.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	// 注册成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  registerMsg,
	})
}

// 用户登录
func UserLogin(c *gin.Context) {
	//获取用户登录信息
	var u model.UserInfo
	u.Telephone = c.PostForm("telephone")
	passwordTemp, ok := c.Get("password")
	u.Password = passwordTemp.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "密码解析错误",
		})
		return
	}

	//获取用户流信息
	loginMsg, err := logic.LoginUser(u.Telephone, u.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	//用户登录成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  loginMsg,
	})
}
