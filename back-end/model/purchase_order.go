package model

// PurchaseOrder Model
type PurchaseOrder struct {
	PurchaseId          int32  `json:"purchase_id"`           // 进货单ID
	UserId              int32  `json:"user_id"`               // 客户ID
	PurchaseDate        string `json:"purchase_date"`         // 进货日期
	PurchaseTotalAmount string `json:"purchase_total_amount"` // 进货总金额
	Node                string `json:"note"`                  // 备注
}

// 指定PurchaseOrder结构体迁移表purchase_order
func (p *PurchaseOrder) TableName() string {
	return "purchase_order"
}
