package controller

import (
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取销售列表
func GetSaleList(c *gin.Context) {
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
	saleList := model.LikeGetSaleListByUserName(search)
	smap := make(map[int32]struct{}) // 去重
	for _, sale := range saleList {
		if _, ok := smap[sale.SaleId]; !ok {
			smap[sale.SaleId] = struct{}{}
		}
	}
	if userId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if sList := model.GetSaleListByUserId(int32(userId)); sList != nil {
			for _, s := range sList {
				if _, ok := smap[s.SaleId]; !ok {
					smap[int32(s.SaleId)] = struct{}{}
					saleList = append(saleList, s)
				}
			}
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, saleList)
}
