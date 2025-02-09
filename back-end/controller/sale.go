package controller

import (
	"drugims/logic"
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

// 获取用户购买记录列表
func GetUserSaleList(c *gin.Context) {
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
	search := fmt.Sprintf("%v", searchValue.(map[string]interface{})["searchValue"])
	saleList := model.LikeGetSaleListByDrugName(search)
	smap := make(map[int32]struct{}) // 去重
	for _, sale := range saleList {
		if _, ok := smap[sale.SaleId]; !ok {
			smap[sale.SaleId] = struct{}{}
		}
	}
	if saleId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if s := model.GetSaleBySaleId(int32(saleId)); s != nil {
			if _, ok := smap[s.SaleId]; !ok {
				smap[int32(s.SaleId)] = struct{}{}
				saleList = append(saleList, s)
			}
		}
	}
	// 过滤userid
	saleListByUserId := []*model.SaleInfo{}
	for _, s := range saleList {
		if s.UserId == userId {
			saleListByUserId = append(saleListByUserId, s)
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, saleList)
}

// 购买药品
func BuyDrug(c *gin.Context) {
	// 获取前端药品信息
	var d model.DrugInfo
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 购买药品
	_, err := logic.BuyDrug(&d)
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
			"msg": "购买成功",
		},
	})
}
