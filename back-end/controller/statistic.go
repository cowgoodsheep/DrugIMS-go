package controller

import (
	"drugims/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取统计信息
func GetStatistics(c *gin.Context) {
	// 获取查询字段
	// var searchValue interface{}
	// if err := c.ShouldBindJSON(&searchValue); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"code": 400,
	// 		"data": gin.H{
	// 			"msg": err.Error(),
	// 		},
	// 	})
	// 	return
	// }
	statistics, err := logic.GetStatistics("1", "2")
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
