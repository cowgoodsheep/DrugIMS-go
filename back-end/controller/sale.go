package controller

import (
	"drugims/logic"
	"drugims/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 创建订单列表
func CreateOrder(c *gin.Context) {
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
	dflow, err := logic.CreateOrder(&d)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 创建成功
	c.JSON(http.StatusOK, dflow.OrderInfo)
}

// 获取订单列表
func GetOrderList(c *gin.Context) {
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
	orderList := model.LikeGetOrderListByUserName(search)
	smap := make(map[int32]struct{}) // 去重
	for _, order := range orderList {
		if _, ok := smap[order.OrderId]; !ok {
			smap[order.OrderId] = struct{}{}
		}
	}
	if userId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if sList := model.GetOrderListByUserId(int32(userId)); sList != nil {
			for _, s := range sList {
				if _, ok := smap[s.OrderId]; !ok {
					smap[int32(s.OrderId)] = struct{}{}
					orderList = append(orderList, s)
				}
			}
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, orderList)
}

// 获取用户购买记录列表
func GetUserOrderList(c *gin.Context) {
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
	orderList := model.LikeGetOrderListByDrugName(search)
	smap := make(map[int32]struct{}) // 去重
	for _, order := range orderList {
		if _, ok := smap[order.OrderId]; !ok {
			smap[order.OrderId] = struct{}{}
		}
	}
	if OrderId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if s := model.GetSaleByOrderId(int32(OrderId)); s != nil {
			if _, ok := smap[s.OrderId]; !ok {
				smap[int32(s.OrderId)] = struct{}{}
				orderList = append(orderList, s)
			}
		}
	}
	// 过滤userid
	dataList := []map[string]interface{}{}
	for _, o := range orderList {
		if o.UserId == userId {
			// 检查订单状态，若为待确认状态则72小时后自动确认，并更新数据库
			if o.OrderStatus == 3 && time.Now().Sub(o.CreateTime).Hours() > 72 {
				_, err := logic.ConfirmOrder(o)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"data": gin.H{
							"msg": err.Error(),
						},
					})
					return
				}
			}

			drug := model.GetDrugByDrugId(o.DrugId)
			dataList = append(dataList, map[string]interface{}{
				"order_id":          o.OrderId,
				"user_id":           o.UserId,
				"user_name":         o.UserName,
				"sale_quantity":     o.SaleQuantity,
				"sale_amount":       o.SaleAmount,
				"supply_amount":     o.SupplyAmount,
				"stock_info":        o.StockInfo,
				"order_status":      o.OrderStatus,
				"create_time":       o.CreateTime,
				"update_time":       o.UpdateTime,
				"drug_id":           drug.DrugId,
				"drug_name":         drug.DrugName,
				"manufacturer":      drug.Manufacturer,
				"unit":              drug.Unit,
				"specification":     drug.Specification,
				"stock_lower_limit": drug.StockLowerLimit,
				"stock_upper_limit": drug.StockUpperLimit,
				"sale_price":        drug.SalePrice,
				"drug_description":  drug.DrugDescription,
				"img":               drug.Img,
				"stock_remain":      drug.StockRemain,
			})
		}
	}
	// 返回数据
	c.JSON(http.StatusOK, dataList)
}

// 确认订单
func ConfirmOrder(c *gin.Context) {
	var o model.OrderInfo
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	_, err := logic.ConfirmOrder(&o)
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
			"msg": "确认成功",
		},
	})
}

// 撤销订单
func RevokeOrder(c *gin.Context) {
	var o model.OrderInfo
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	_, err := logic.RevokeOrder(&o)
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
			"msg": "撤销成功",
		},
	})
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

// 退款
func RefundOrder(c *gin.Context) {
	var o model.OrderInfo
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	_, err := logic.RefundOrder(&o)
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
			"msg": "退款申请发起成功",
		},
	})
}
