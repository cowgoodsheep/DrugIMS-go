package controller

import (
	"drugims/logic"
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 供应药品
func SupplyDrug(c *gin.Context) {
	// 获取供应信息
	var s model.SupplyOrder
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	_, err := logic.SupplyDrug(&s)
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
			"msg": "供应成功",
		},
	})
}

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

// 获取进货列表
func GetSupplyList(c *gin.Context) {
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
	saleList := model.LikeGetSupplyListByUserName(search)
	// 返回数据
	c.JSON(http.StatusOK, saleList)
}

// 获取用户进货记录列表
func GetUserSupplyList(c *gin.Context) {
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
	userId := int32(searchValue.(map[string]interface{})["user_id"].(float64))
	startDate := fmt.Sprintf("%v", searchValue.(map[string]interface{})["startDate"])
	endDate := fmt.Sprintf("%v", searchValue.(map[string]interface{})["endDate"])
	supplyList := model.GetSupplyListByTime(startDate, endDate)
	smap := make(map[int32]struct{}) // 去重
	for _, supply := range supplyList {
		if _, ok := smap[supply.SupplyId]; !ok {
			smap[supply.SupplyId] = struct{}{}
		}
	}
	// 过滤userid
	supplyListByUserId := []*model.SupplyOrder{}
	for _, s := range supplyList {
		if s.UserId == userId {
			supplyListByUserId = append(supplyListByUserId, s)
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, supplyListByUserId)
}
