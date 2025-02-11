package model

import (
	"drugims/dao"
	"errors"
	"time"
)

// SaleInfo Model
type SaleInfo struct {
	SaleId       int32     `json:"sale_id"`       // 销售ID
	DrugId       int32     `json:"drug_id"`       // 药品ID
	UserId       int32     `json:"user_id"`       // 客户ID
	SaleQuantity int32     `json:"sale_quantity"` // 销售数量
	SaleAmount   float32   `json:"sale_amount"`   // 销售金额
	SupplyAmount float32   `json:"supply_amount"` // 进货金额
	CreateTime   time.Time `json:"create_time" gorm:"-"`
	UpdateTime   time.Time `json:"update_time" gorm:"-"`

	UserName string `json:"user_name" gorm:"-"` // 用户名称
	DrugName string `json:"drug_name" gorm:"-"` // 药品名称
}

// 指定SaleInfo结构体迁移表sale_info
func (s *SaleInfo) TableName() string {
	return "sale_info"
}

// GetSaleBySaleId 获取根据销售id获取销售记录
func GetSaleBySaleId(saleId int32) *SaleInfo {
	var sFind SaleInfo
	dao.DB.Model(&SaleInfo{}).Where("sale_id=?", saleId).First(&sFind)
	if sFind.SaleId == 0 {
		return nil
	}
	dFind := GetDrugByDrugId(sFind.DrugId)
	sFind.DrugName = dFind.DrugName
	return &sFind
}

// GetSaleListByUserId 获取根据用户id获取销售信息
func GetSaleListByUserId(userId int32) []*SaleInfo {
	var sListFind []*SaleInfo
	dao.DB.Model(&SaleInfo{}).Where("user_id=?", userId).Find(&sListFind)
	uFind := GetUserByUserId(userId)
	for _, s := range sListFind {
		s.UserName = uFind.UserName
		dFind := GetDrugByDrugId(s.DrugId)
		s.DrugName = dFind.DrugName
	}
	return sListFind
}

// GetSaleListByDrugId 获取根据药品id获取销售信息
func GetSaleListByDrugId(drugId int32) []*SaleInfo {
	var sListFind []*SaleInfo
	dao.DB.Model(&SaleInfo{}).Where("drug_id=?", drugId).Find(&sListFind)
	for _, s := range sListFind {
		dFind := GetDrugByDrugId(s.DrugId)
		uFind := GetUserByUserId(s.UserId)
		s.DrugName = dFind.DrugName
		s.UserName = uFind.UserName
	}
	return sListFind
}

// LikeGetSaleListByUserName 获取根据用户名称模模糊查询销售列表
func LikeGetSaleListByUserName(userName string) []*SaleInfo {
	// 获取用户id
	uListFind := LikeGetUserListByUserName(userName)
	sListFind := []*SaleInfo{}
	for _, u := range uListFind {
		sList := GetSaleListByUserId(u.UserId)
		if len(sList) > 0 {
			sListFind = append(sListFind, sList...)
		}
	}
	return sListFind
}

// LikeGetSaleListByDrugName 获取根据药品名称模模糊查询销售列表
func LikeGetSaleListByDrugName(drugName string) []*SaleInfo {
	// 获取药品id
	dListFind := LikeGetDrugListByDrugName(drugName)
	sListFind := []*SaleInfo{}
	for _, d := range dListFind {
		sList := GetSaleListByDrugId(d.DrugId)
		if len(sList) > 0 {
			sListFind = append(sListFind, sList...)
		}
	}
	return sListFind
}

// CreateSale 创建销售记录
func CreateSale(s *SaleInfo) error {
	if s == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Model(&SaleInfo{}).Create(s).Error
}

// GetSaleListByTime 获取根据日期查询销售列表
func GetSaleListByTime(startDate string, endDate string) []*SaleInfo {
	sListFind := []*SaleInfo{}
	if startDate != "" && endDate != "" {
		dao.DB.Where("create_time BETWEEN ? AND ?", startDate, endDate).Find(&sListFind)
	} else {
		dao.DB.Find(&sListFind)
	}
	for _, s := range sListFind {
		uFind := GetUserByUserId(s.UserId)
		s.UserName = uFind.UserName
		dFind := GetDrugByDrugId(s.DrugId)
		s.DrugName = dFind.DrugName
	}
	return sListFind
}
