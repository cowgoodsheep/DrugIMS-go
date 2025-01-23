package logic

import (
	"drugims/model"
)

// DrugInfoFlow 药品信息流
type DrugInfoFlow struct {
	DrugInfo *model.DrugInfo
}

// 添加药品
func AddDrug(drugInfo *model.DrugInfo) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(drugInfo).addDrug()
}

// 修改药品
func UpdateDrug(drugInfo *model.DrugInfo) (*DrugInfoFlow, error) {
	return NewDrugInfoFlow(drugInfo).updateDrug()
}

func NewDrugInfoFlow(drugInfo *model.DrugInfo) *DrugInfoFlow {
	return &DrugInfoFlow{DrugInfo: drugInfo}
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
		Price:           d.DrugInfo.Price,
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
		Price:           d.DrugInfo.Price,
		DrugDescription: d.DrugInfo.DrugDescription,
	}
	if err := model.UpdateDrug(d.DrugInfo.DrugId, drugInfo); err != nil {
		return d, err
	}
	return d, nil
}
