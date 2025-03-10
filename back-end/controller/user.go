package controller

import (
	"drugims/logic"
	"drugims/middleware"
	"drugims/model"
	"fmt"
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
	registerMsg, err := logic.RegisterUser(&u)
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
	// 获取用户登录信息
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
	// 获取用户流信息
	loginMsg, err := logic.LoginUser(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 更新用户token
	u.Token = loginMsg.Token
	model.UpdateUserInfo(u.UserId, &u)
	c.Header("Token", loginMsg.Token)
	c.JSON(http.StatusOK, []gin.H{{
		"address":   loginMsg.UserInfo.Address,
		"telephone": loginMsg.UserInfo.Telephone,
		"password":  loginMsg.UserInfo.Password,
		"role":      loginMsg.UserInfo.Role,
		"user_id":   loginMsg.UserInfo.UserId,
		"user_name": loginMsg.UserInfo.UserName,
	}})
}

// 获取个人信息
func GetUser(c *gin.Context) {
	// 获取用户信息
	var userId int32
	if err := c.ShouldBindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, model.QueryUserByUserId(userId))
}

// 用户信息更新
func UpdateUser(c *gin.Context) {
	// 获取用户信息
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
	// 获取用户流信息
	_, err := logic.UpdateUser(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}

// 获取用户列表
func GetUserList(c *gin.Context) {
	// 获取查询字段
	var searchValue interface{}
	if err := c.ShouldBindJSON(&searchValue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	search := fmt.Sprintf("%v", searchValue)
	userList := model.LikeGetUserListByUserName(search)
	c.JSON(http.StatusOK, userList)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	// 获取用户信息
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
	model.DeleteUser(u.UserId)
	c.JSON(http.StatusOK, nil)
}

// 拉黑用户
func BlockUser(c *gin.Context) {
	// 获取用户信息
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
	// 更新用户账户状态
	user := model.QueryUserByUserId(u.UserId)
	user.Status = 2
	model.UpdateUserInfo(user.UserId, user)
	// 将旧token暂时存入内存的黑名单，强制踢出登录
	middleware.Blacklist.Store(user.Token, true)
	c.JSON(http.StatusOK, nil)
}

// 解除拉黑
func UnblockUser(c *gin.Context) {
	// 获取用户信息
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
	// 更新用户账户状态
	user := model.QueryUserByUserId(u.UserId)
	user.Status = 1
	model.UpdateUserInfo(user.UserId, user)
	c.JSON(http.StatusOK, nil)
}

// 充值
func Recharge(c *gin.Context) {
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
	_, err := logic.Recharge(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg": "充值成功",
		},
	})
}

// 提现
func Withdraw(c *gin.Context) {
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
	_, err := logic.Withdraw(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg": "提现成功",
		},
	})
}
