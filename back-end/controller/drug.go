package controller

import (
	"drugims/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取药品列表
func GetDrugList(c *gin.Context) {
	drugList := model.GetDrugList()
	c.JSON(http.StatusOK, drugList)
}
