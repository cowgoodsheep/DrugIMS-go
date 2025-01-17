package model

import "drugims/dao"

// StockInfo Model
type StockInfo struct {
	StockId           int32  `json:"stock_id"`            // 库存ID
	DrugId            int32  `json:"drug_id"`             // 药品ID
	BatchNumber       string `json:"batch_number"`        // 批号
	ProductionDate    string `json:"production_date"`     // 生产日期
	PurchaseDate      string `json:"purchase_date"`       // 进货日期
	PurchaseUnitPrice string `json:"purchase_unit_price"` // 进货单价
	RemainingQuantity int32  `json:"remaining_quantity"`  // 剩余数量
}

// 指定StockInfo结构体迁移表stock_info
func (s *StockInfo) TableName() string {
	return "stock_info"
}

// GetDrugRemain 获取药品剩余数量
func GetDrugRemain(drugId int32) int64 {
	var total int64
	row := dao.DB.Model(&StockInfo{}).Where("drug_id = ?", drugId).Select("SUM(remaining_quantity) as total").Row()
	row.Scan(&total)
	return total
}
