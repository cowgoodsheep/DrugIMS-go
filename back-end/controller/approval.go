package controller

import (
	"drugims/logic"
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取审批列表
func GetApprovalList(c *gin.Context) {
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
	approvalList := []*model.ApprovalInfo{}

	// 根据用户姓名模糊搜索用户id列表
	userList := model.LikeGetUserListByUserName(search)
	amap := make(map[int32]struct{}) // 去重
	for _, user := range userList {
		amap[user.UserId] = struct{}{}
		approvalList = append(approvalList, model.GetApprovalListByApprovalId(user.UserId)...)
	}
	if userId, err := strconv.Atoi(search); err == nil { // 为id继续搜索id
		if _, ok := amap[int32(userId)]; !ok {
			approvalList = append(approvalList, model.GetApprovalListByApprovalId(int32(userId))...)
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, approvalList)
}

// 审批操作
func ApprovalOperate(c *gin.Context) {
	var a model.ApprovalInfo
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	err := logic.ApprovalOperate(&a)
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
			"msg": "操作成功",
		},
	})
}
