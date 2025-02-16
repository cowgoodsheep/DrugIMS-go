package logic

import (
	"drugims/model"
	"fmt"
	"sort"
	"time"
)

// DrugInfoFlow 药品信息流
type DrugInfoFlow struct {
	DrugInfo    *model.DrugInfo
	SupplyOrder *model.SupplyOrder
}

// 添加药品
func AddDrug(drugInfo *model.DrugInfo) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(drugInfo, nil).addDrug()
}

// 修改药品
func UpdateDrug(drugInfo *model.DrugInfo) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(drugInfo, nil).updateDrug()
}

// 购买药品
func BuyDrug(drugInfo *model.DrugInfo) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(drugInfo, nil).buyDrug()
}

// 供应药品
func SupplyDrug(supplyOrder *model.SupplyOrder) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(nil, supplyOrder).supplyDrug()
}

func NewDrugInfoFlow(drugInfo *model.DrugInfo, supplyOrder *model.SupplyOrder) *DrugInfoFlow {
	return &DrugInfoFlow{DrugInfo: drugInfo, SupplyOrder: supplyOrder}
}

// 添加
func (d *DrugInfoFlow) addDrug() (*DrugInfoFlow, error) {
	drugInfo := &model.DrugInfo{
		DrugName:        d.DrugInfo.DrugName,
		Manufacturer:    d.DrugInfo.Manufacturer,
		Unit:            d.DrugInfo.Unit,
		Specification:   d.DrugInfo.Specification,
		StockLowerLimit: d.DrugInfo.StockLowerLimit,
		StockUpperLimit: d.DrugInfo.StockUpperLimit,
		SalePrice:       d.DrugInfo.SalePrice,
		DrugDescription: d.DrugInfo.DrugDescription,
	}
	if err := model.CreateDrug(drugInfo); err != nil {
		return d, err
	}
	return d, nil
}

// 修改
func (d *DrugInfoFlow) updateDrug() (*DrugInfoFlow, error) {
	drugInfo := &model.DrugInfo{
		DrugName:        d.DrugInfo.DrugName,
		Manufacturer:    d.DrugInfo.Manufacturer,
		Unit:            d.DrugInfo.Unit,
		Specification:   d.DrugInfo.Specification,
		StockLowerLimit: d.DrugInfo.StockLowerLimit,
		StockUpperLimit: d.DrugInfo.StockUpperLimit,
		SalePrice:       d.DrugInfo.SalePrice,
		DrugDescription: d.DrugInfo.DrugDescription,
	}
	if err := model.UpdateDrug(d.DrugInfo.DrugId, drugInfo); err != nil {
		return d, err
	}
	return d, nil
}

// 购买
func (d *DrugInfoFlow) buyDrug() (*DrugInfoFlow, error) {
	// 检查剩余数量是否足够
	if d.DrugInfo.StockRemain < d.DrugInfo.SaleQuantity {
		return nil, fmt.Errorf("库存不足")
	}
	// 取出库存
	stockList := model.GetStockListByDrugId(d.DrugInfo.DrugId)
	// 扣之前再检查一遍
	totalStock := int32(0)
	for _, stock := range stockList {
		totalStock += stock.RemainingQuantity
	}
	if totalStock < d.DrugInfo.SaleQuantity {
		return nil, fmt.Errorf("库存不足")
	}
	// 根据生产日期从小到大排序
	sort.Slice(stockList, func(i, j int) bool {
		dateI, err1 := time.Parse("2006-01-02", stockList[i].ProductionDate)
		dateJ, err2 := time.Parse("2006-01-02", stockList[j].ProductionDate)
		if err1 != nil || err2 != nil {
			return stockList[i].StockId > stockList[j].StockId
		}
		return dateI.Before(dateJ)
	})
	// 从旧到新，一个个更新库存数量
	buyNum := d.DrugInfo.SaleQuantity // 购买数量
	saleTotalPrice := float32(0)      // 销售总价
	supplyTotalPrice := float32(0)    // 采购总价
	for _, stock := range stockList {
		if stock.RemainingQuantity > buyNum {
			stock.RemainingQuantity -= buyNum
			saleTotalPrice += float32(buyNum) * d.DrugInfo.SalePrice
			supplyTotalPrice += float32(buyNum) * stock.SupplyPrice
			// 更新改库存的数量
			if err := model.UpdateStock(stock); err != nil {
				return nil, err
			}
			break
		} else {
			buyNum -= stock.RemainingQuantity
			saleTotalPrice += float32(stock.RemainingQuantity) * d.DrugInfo.SalePrice
			supplyTotalPrice += float32(stock.RemainingQuantity) * stock.SupplyPrice
			// 直接删掉改库存数据
			if err := model.DeleteStock(stock.StockId); err != nil {
				return nil, err
			}
		}
	}
	// 生成销售记录
	saleInfo := &model.SaleInfo{
		DrugId:       d.DrugInfo.DrugId,
		UserId:       d.DrugInfo.UserId,
		SaleQuantity: d.DrugInfo.SaleQuantity,
		SaleAmount:   saleTotalPrice,
		SupplyAmount: supplyTotalPrice,
	}
	if err := model.CreateSale(saleInfo); err != nil {
		return nil, err
	}
	return d, nil
}

// 供应
func (d *DrugInfoFlow) supplyDrug() (*DrugInfoFlow, error) {
	d.DrugInfo = model.GetDrugByDrugId(d.SupplyOrder.DrugId)
	// 检查供应价格是否高于售价
	if d.DrugInfo.SalePrice <= d.SupplyOrder.SupplyPrice {
		return nil, fmt.Errorf("供应价格不得高于售价")
	}
	// 检查供应数量加上库存剩余是否超出库存上限
	drugRemain := model.GetDrugRemain(d.SupplyOrder.DrugId)
	if drugRemain+d.SupplyOrder.SupplyQuantity > d.DrugInfo.StockUpperLimit {
		return nil, fmt.Errorf("不得超出库存上限")
	}
	// 时间格式化为指定格式
	t, err := time.Parse("2006-01-02T15:04:05.999Z", d.SupplyOrder.ProductionDate)
	if err != nil {
		return nil, err
	}
	d.SupplyOrder.ProductionDate = t.Format("2006-01-02")
	// 存储库存信息
	if err := model.CreateStock(&model.StockInfo{
		DrugId:            d.DrugInfo.DrugId,
		BatchNumber:       d.SupplyOrder.BatchNumber,
		ProductionDate:    d.SupplyOrder.ProductionDate,
		SupplyPrice:       d.SupplyOrder.SupplyPrice,
		RemainingQuantity: d.SupplyOrder.SupplyQuantity,
	}); err != nil {
		return nil, err
	}
	// 存储供应记录
	if err := model.CreateSupply(d.SupplyOrder); err != nil {
		return nil, err
	}
	return d, nil
}
