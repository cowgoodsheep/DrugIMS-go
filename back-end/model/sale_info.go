package model

import (
	"drugims/dao"
)

// SaleInfo Model
type SaleInfo struct {
	SaleId        int32  `json:"sale_id"`         // 销售ID
	DrugId        int32  `json:"drug_id"`         // 药品ID
	UserId        int32  `json:"user_id"`         // 客户ID
	SaleDate      string `json:"sale_date"`       // 销售日期
	SaleQuantity  int32  `json:"sale_quantity"`   // 销售数量
	SaleUnitPrice string `json:"sale_unit_price"` // 销售单价
	SaleAmount    string `json:"sale_amount"`     // 销售金额

	UserName string `json:"user_name" gorm:"-"` // 用户名称
}

// 指定SaleInfo结构体迁移表sale_info
func (s *SaleInfo) TableName() string {
	return "sale_info"
}

// GetSaleListByUserId 获取根据用户id获取销售信息
func GetSaleListByUserId(userId int32) []*SaleInfo {
	var sFind []*SaleInfo
	dao.DB.Model(&SaleInfo{}).Where("user_id=?", userId).Find(&sFind)
	return sFind
}

// LikeGetSaleListByUserName 获取根据用户名称模模糊查询销售列表
func LikeGetSaleListByUserName(userName string) []*SaleInfo {
	// 先获取用户id
	uListFind := LikeGetUserListByUserName(userName)
	sListFind := []*SaleInfo{}
	for _, u := range uListFind {
		sList := GetSaleListByUserId(u.UserId)
		for _, s := range sList {
			s.UserName = u.UserName
		}
		if len(sList) > 0 {
			sListFind = append(sListFind, sList...)
		}
	}
	return sListFind
}
