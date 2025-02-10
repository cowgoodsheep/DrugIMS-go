package model

import (
	"drugims/dao"
	"errors"
)

// DrugInfo Model
type DrugInfo struct {
	DrugId          int32   `json:"drug_id"`           // 药品ID
	DrugName        string  `json:"drug_name"`         // 药品名称
	Manufacturer    string  `json:"manufacturer"`      // 生产厂家
	Unit            string  `json:"unit"`              // 单位
	Specification   string  `json:"specification"`     // 规格
	StockLowerLimit int32   `json:"stock_lower_limit"` // 库存下限
	StockUpperLimit int32   `json:"stock_upper_limit"` // 库存上限
	SalePrice       float32 `json:"sale_price"`        // 售价
	DrugDescription string  `json:"drug_description"`  // 药品说明
	Img             string  `json:"img"`               // 药品图片
	CreateTime      string  `json:"create_time"`
	UpdateTime      string  `json:"update_time"`

	StockRemain  int32 `json:"stock_remain" gorm:"-"`  // 库存剩余数量
	SaleQuantity int32 `json:"sale_quantity" gorm:"-"` // 购买数量
	UserId       int32 `json:"user_id" gorm:"-"`       // 购买用户
}

// 指定DrugInfo结构体迁移表user
func (d *DrugInfo) TableName() string {
	return "drug_info"
}

// GetDrugByDrugId 获取根据药品id获取药品
func GetDrugByDrugId(drugId int32) *DrugInfo {
	var dFind DrugInfo
	dao.DB.Model(&DrugInfo{}).Where("drug_id=?", drugId).First(&dFind)
	if dFind.DrugId == 0 {
		return nil
	}
	return &dFind
}

// LikeGetDrugListByDrugName 获取根据药品名称模模糊查询药品列表
func LikeGetDrugListByDrugName(drugName string) []*DrugInfo {
	var dListFind []*DrugInfo
	dao.DB.Model(&DrugInfo{}).Where("drug_name LIKE ?", "%"+drugName+"%").Find(&dListFind)
	return dListFind
}

// CreateDrug 创建药品
func CreateDrug(d *DrugInfo) error {
	if d == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Model(&DrugInfo{}).Create(d).Error
}

// UpdateDrug 更新药品数据
func UpdateDrug(drugId int32, d *DrugInfo) error {
	if d == nil {
		return errors.New("空指针错误")
	}
	dao.DB.Model(&DrugInfo{}).Where("drug_id = ?", drugId).Updates(d)
	return nil
}

// DeleteDrug 删除药品
func DeleteDrug(drugId int32) error {
	return dao.DB.Model(&DrugInfo{}).Where("drug_id = ?", drugId).Delete(&DrugInfo{}).Error
}
