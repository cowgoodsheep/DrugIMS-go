package model

import (
	"drugims/dao"
	"errors"
)

// StockInfo Model
type StockInfo struct {
	StockId           int32   `json:"stock_id"`           // 库存ID
	DrugId            int32   `json:"drug_id"`            // 药品ID
	BatchNumber       string  `json:"batch_number"`       // 批号
	ProductionDate    string  `json:"production_date"`    // 生产日期
	SupplyDate        string  `json:"supply_date"`        // 进货日期
	SupplyPrice       float32 `json:"supply_price"`       // 进货单价
	RemainingQuantity int32   `json:"remaining_quantity"` // 剩余数量

	DrugName string `json:"drug_name" gorm:"-"` // 库存剩余数量
}

// 指定StockInfo结构体迁移表stock_info
func (s *StockInfo) TableName() string {
	return "stock_info"
}

// GetDrugRemain 获取药品剩余数量
func GetDrugRemain(drugId int32) int32 {
	var total int32
	row := dao.DB.Model(&StockInfo{}).Where("drug_id = ?", drugId).Select("SUM(remaining_quantity) as total").Row()
	row.Scan(&total)
	return total
}

// GetStockListByDrugId 获取根据药品id获取库存
func GetStockListByDrugId(drugId int32) []*StockInfo {
	var sFind []*StockInfo
	dao.DB.Model(&StockInfo{}).Where("drug_id=?", drugId).Find(&sFind)
	return sFind
}

// LikeGetStockListByDrugName 获取根据药品名称模模糊查询库存列表
func LikeGetStockListByDrugName(drugName string) []*StockInfo {
	// 先获取药品id
	dListFind := LikeGetDrugListByDrugName(drugName)
	sListFind := []*StockInfo{}
	for _, d := range dListFind {
		sList := GetStockListByDrugId(d.DrugId)
		for _, s := range sList {
			s.DrugName = d.DrugName
		}
		if len(sList) > 0 {
			sListFind = append(sListFind, sList...)
		}
	}
	return sListFind
}

// GetStockByStockId 获取根据库存id获取药品
func GetStockByStockId(stockId int32) *StockInfo {
	var sFind StockInfo
	dao.DB.Model(&StockInfo{}).Where("stock_id=?", stockId).First(&sFind)
	if sFind.StockId == 0 {
		return nil
	}
	d := GetDrugByDrugId(sFind.DrugId)
	if d.DrugId == 0 {
		return nil
	}
	sFind.DrugName = d.DrugName
	return &sFind
}

// UpdateStock 更新库存
func UpdateStock(stockInfo *StockInfo) error {
	return dao.DB.Model(&StockInfo{}).Where("stock_id = ?", stockInfo.StockId).Updates(stockInfo).Error
}

// DeleteStock 删除库存
func DeleteStock(stockId int32) error {
	return dao.DB.Model(&StockInfo{}).Where("stock_id = ?", stockId).Delete(&StockInfo{}).Error
}

// CreateStock 创建库存记录
func CreateStock(s *StockInfo) error {
	if s == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Model(&StockInfo{}).Create(s).Error
}
