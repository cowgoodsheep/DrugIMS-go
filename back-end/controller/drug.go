package controller

import (
	"drugims/logic"
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
	drugList := model.LikeGetDrugListByDrugName(search)
	dmap := make(map[int32]struct{}) // 去重
	for _, drug := range drugList {
		if _, ok := dmap[drug.DrugId]; !ok {
			dmap[drug.DrugId] = struct{}{}
		}
	}
	if drugId, err := strconv.Atoi(search); err == nil { // 为id则继续搜索id
		if _, ok := dmap[int32(drugId)]; !ok {
			if d := model.GetDrugByDrugId(int32(drugId)); d != nil {
				dmap[int32(drugId)] = struct{}{}
				drugList = append(drugList, d)
			}
		}
	}
	// 获取药品剩余数量
	for i, drug := range drugList {
		drugList[i].StockRemain = model.GetDrugRemain(drug.DrugId)
	}
	// 返回数据
	c.JSON(http.StatusOK, drugList)
}

// 添加药品
func AddDrug(c *gin.Context) {
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
	// 上传药品信息至添加药品服务
	_, err := logic.AddDrug(&d)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 添加成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg": "药品信息添加成功",
		},
	})
}

// 修改药品
func UpdateDrug(c *gin.Context) {
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
	// 上传药品信息至添加药品服务
	_, err := logic.UpdateDrug(&d)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"data": gin.H{
				"msg": err.Error(),
			},
		})
		return
	}
	// 添加成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"msg": "药品信息修改成功",
		},
	})
}

// 删除药品
func DeleteDrug(c *gin.Context) {
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
	if err := model.DeleteDrug(d.DrugId); err != nil {
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": gin.H{
					"msg": err.Error(),
				},
			})
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}
