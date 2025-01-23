package controller

import (
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取库存列表
func GetStockList(c *gin.Context) {
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
	stockList := model.LikeGetStockListByDrugName(search)
	smap := make(map[int32]struct{}) // 去重
	for _, stock := range stockList {
		if _, ok := smap[stock.StockId]; !ok {
			smap[stock.StockId] = struct{}{}
		}
	}
	if stockId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if _, ok := smap[int32(stockId)]; !ok {
			if s := model.GetStockByStockId(int32(stockId)); s != nil {
				smap[int32(stockId)] = struct{}{}
				stockList = append(stockList, s)
			}
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, stockList)
}
