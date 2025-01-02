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
	Price           int32  `json:"price"`             // 售价
	DrugDescription string `json:"drug_description"`  // 药品说明
	Img             int8   `json:"img"`               // 药品图片
}

// 指定DrugInfo结构体迁移表user
func (u *DrugInfo) TableName() string {
	return "drug_info"
}

// GetDrugList 获取药品列表
func GetDrugList() *[]DrugInfo {
	var dListFind []DrugInfo
	dao.DB.Find(&dListFind)
	return &dListFind
}
