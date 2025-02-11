package model

import (
	"drugims/dao"
	"errors"
	"time"
)

// SupplyOrder Model
type SupplyOrder struct {
	SupplyId       int32     `json:"supply_id"`       // 进货单ID
	DrugId         int32     `json:"drug_id"`         // 药品ID
	UserId         int32     `json:"user_id"`         // 客户ID
	BatchNumber    string    `json:"batch_number"`    // 批号
	ProductionDate string    `json:"production_date"` // 生产日期
	SupplyQuantity int32     `json:"supply_quantity"` // 进货数量
	SupplyPrice    float32   `json:"supply_price"`    // 进货单价
	Note           string    `json:"note"`            // 备注
	CreateTime     time.Time `json:"create_time" gorm:"-"`
	UpdateTime     time.Time `json:"update_time" gorm:"-"`

	UserName string `json:"user_name" gorm:"-"` // 用户名称
	DrugName string `json:"drug_name" gorm:"-"` // 药品名称
}

// 指定SupplyOrder结构体迁移表supply_order
func (s *SupplyOrder) TableName() string {
	return "supply_order"
}

// CreateSupply 创建供应记录
func CreateSupply(s *SupplyOrder) error {
	if s == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Model(&SupplyOrder{}).Create(s).Error
}

// GetSupplyListByUserId 获取根据用户id获取进货信息
func GetSupplyListByUserId(userId int32) []*SupplyOrder {
	var sListFind []*SupplyOrder
	dao.DB.Model(&SupplyOrder{}).Where("user_id=?", userId).Find(&sListFind)
	uFind := GetUserByUserId(userId)
	for _, s := range sListFind {
		s.UserName = uFind.UserName
		dFind := GetDrugByDrugId(s.DrugId)
		s.DrugName = dFind.DrugName
	}
	return sListFind
}

// LikeGetSupplyListByUserName 获取根据用户名称模糊查询进货列表
func LikeGetSupplyListByUserName(userName string) []*SupplyOrder {
	// 获取用户id
	uListFind := LikeGetUserListByUserName(userName)
	sListFind := []*SupplyOrder{}
	for _, u := range uListFind {
		sList := GetSupplyListByUserId(u.UserId)
		if len(sList) > 0 {
			sListFind = append(sListFind, sList...)
		}
	}
	return sListFind
}

// GetSupplyListByTime 获取根据日期查询进货列表
func GetSupplyListByTime(startDate string, endDate string) []*SupplyOrder {
	sListFind := []*SupplyOrder{}
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
