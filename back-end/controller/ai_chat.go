package controller

import (
	"drugims/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AiChat(c *gin.Context) {
	// 获取提问信息
	var q []logic.Message
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	response, err := logic.GetAiChatResponse(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}

	// 返回最新回复
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg": response,
		},
	})
}
