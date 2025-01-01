package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 主页显示
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "你好",
	})
}
