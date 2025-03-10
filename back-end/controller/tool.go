package controller

import (
	"drugims/logic"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取统计信息
func GetStatistics(c *gin.Context) {
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
	startDate := fmt.Sprintf("%v", searchValue.(map[string]interface{})["startDate"])
	endDate := fmt.Sprintf("%v", searchValue.(map[string]interface{})["endDate"])
	statistics, err := logic.GetStatisticByTime(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"statistics": statistics,
		},
	})
}

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

// RiskManage
func RiskManage(c *gin.Context) {
	c.JSON(http.StatusOK, logic.RiskManage())
}
