package logic

import (
	"drugims/model"
)

// DrugInfoFlow 药品信息流
type DrugInfoFlow struct {
	DrugInfo *model.DrugInfo
	DrugId   int32 `json:"drug_id"`
}
