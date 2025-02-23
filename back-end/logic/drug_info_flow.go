package logic

import (
	"drugims/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
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

// 供应
func (d *DrugInfoFlow) supplyDrug() (*DrugInfoFlow, error) {
	d.DrugInfo = model.GetDrugByDrugId(d.SupplyOrder.DrugId)
	// 检查供应价格是否高于售价
	if d.DrugInfo.SalePrice.LessThanOrEqual(d.SupplyOrder.SupplyPrice) {
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
	d.SupplyOrder.SupplyStatus = 0
	// 存储供应记录
	if dt, err := model.CreateSupply(d.SupplyOrder); err != nil {
		return nil, err
	} else {
		d.SupplyOrder = model.GetSupplyBySupplyId(dt.SupplyId)
		d.SupplyOrder.ProductionDate = t.Format("2006-01-02")
	}

	// 添加冻结资金
	u := model.QueryUserByUserId(d.SupplyOrder.UserId)
	u.BlockBalance = u.BlockBalance.Add(d.SupplyOrder.SupplyPrice.Mul(decimal.NewFromInt(int64(d.SupplyOrder.SupplyQuantity))))
	model.UpdateUserInfo(u.UserId, u)

	// 发起进货审批
	d.SupplyOrder.DrugName = d.DrugInfo.DrugName
	d.SupplyOrder.UserName = u.UserName
	supplyOrderJson, err := json.Marshal(d.SupplyOrder)
	if err != nil {
		return nil, err
	}
	approvalInfo := &model.ApprovalInfo{
		UserId:          d.SupplyOrder.UserId,
		ApprovalType:    1,
		ApprovalContent: string(supplyOrderJson),
		Reason:          d.SupplyOrder.Note,
		ApprovalStatus:  0,
	}
	model.CreateApproval(approvalInfo)
	return d, nil
}
