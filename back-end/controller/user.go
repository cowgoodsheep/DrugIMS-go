package controller

import (
	"drugims/logic"
	"drugims/middleware"
	"drugims/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(c *gin.Context) {
	// 获取用户注册信息
	var u model.UserInfo
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 对密码进行加密
	if password, err := middleware.SHAMiddleWare(u.Password); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	} else {
		u.Password = password
	}
	// 上传用户注册信息至用户登录服务，进行用户注册
	registerMsg, err := logic.RegisterUser(u.Telephone, u.UserName, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 注册成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg":         "注册成功",
			"registerMsg": registerMsg,
		},
	})
}

// 用户登录
func UserLogin(c *gin.Context) {
	//获取用户登录信息
	var u model.UserInfo
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 对密码进行加密
	if password, err := middleware.SHAMiddleWare(u.Password); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	} else {
		u.Password = password
	}
	//获取用户流信息
	loginMsg, err := logic.LoginUser(u.Telephone, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	//用户登录成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"loginMsg": loginMsg,
		},
	})
}
