package controller

import (
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取药品列表
func GetDrugList(c *gin.Context) {
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
	// 判断查询字段是id还是名称
	drugList := &[]model.DrugInfo{}
	if drugId, err := strconv.Atoi(search); err != nil { // 不为id则模糊搜索名称
		drugList = model.LikeGetDrugListByDrugName(search)
	} else {
		*drugList = append(*drugList, *model.GetDrugByDrugId(int32(drugId)))
	}
	// 获取药品剩余数量
	for i, drug := range *drugList {
		(*drugList)[i].StockRemain = model.GetDrugRemain(drug.DrugId)
	}
	// 返回数据
	c.JSON(http.StatusOK, drugList)
}
