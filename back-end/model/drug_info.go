package model

import "drugims/dao"

// DrugInfo Model
type DrugInfo struct {
	DrugId          int32  `json:"drug_id"`           // 药品ID
	DrugName        string `json:"drug_name"`         // 药品名称
	Manufacturer    string `json:"manufacturer"`      // 生产厂家
	Unit            string `json:"unit"`              // 单位
	Specification   string `json:"specification"`     // 规格
	StockLowerLimit int32  `json:"stock_lower_limit"` // 库存下限
	StockUpperLimit int32  `json:"stock_upper_limit"` // 库存上限
	StockRemain     int64  `json:"stock_remain"`      // 库存剩余数量
	Price           int32  `json:"price"`             // 售价
	DrugDescription string `json:"drug_description"`  // 药品说明
	Img             string `json:"img"`               // 药品图片
}

// 指定DrugInfo结构体迁移表user
func (d *DrugInfo) TableName() string {
	return "drug_info"
}

// GetDrugByDrugId 获取根据药品id获取药品
func GetDrugByDrugId(drugId int32) *DrugInfo {
	var dFind DrugInfo
	dao.DB.Model(&DrugInfo{}).Where("drug_id=?", drugId).First(&dFind)
	return &dFind
}

// LikeGetDrugListLike 获取根据药品名称模模糊查询药品列表
func LikeGetDrugListByDrugName(drugName string) *[]DrugInfo {
	var dListFind []DrugInfo
	dao.DB.Model(&DrugInfo{}).Where("drug_name LIKE ?", "%"+drugName+"%").Find(&dListFind)
	return &dListFind
}
