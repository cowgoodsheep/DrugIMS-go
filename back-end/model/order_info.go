package model

import (
	"drugims/dao"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// OrderInfo Model
type OrderInfo struct {
	OrderId      int32           `json:"order_id" gorm:"primary_key;auto_increment"`                   // 订单ID
	DrugId       int32           `json:"drug_id"`                                                      // 药品ID
	UserId       int32           `json:"user_id"`                                                      // 客户ID
	SaleQuantity int32           `json:"sale_quantity"`                                                // 订单数量
	SaleAmount   decimal.Decimal `json:"sale_amount" gorm:"type:decimal(10,2);column:sale_amount"`     // 订单金额
	SupplyAmount decimal.Decimal `json:"supply_amount" gorm:"type:decimal(10,2);column:supply_amount"` // 进货金额
	OrderStatus  int32           `json:"order_status"`                                                 // 订单状态,0处理中,1已完成,2已撤销,3待确认
	StockInfo    string          `json:"stock_info"`                                                   // 扣除库存
	CreateTime   time.Time       `json:"create_time" gorm:"-"`
	UpdateTime   time.Time       `json:"update_time" gorm:"-"`

	UserName string `json:"user_name" gorm:"-"` // 用户名称
	DrugName string `json:"drug_name" gorm:"-"` // 药品名称
	Reason   string `json:"reason" gorm:"-"`    // 退款理由
}

// 指定OrderInfo结构体迁移表sale_info
func (s *OrderInfo) TableName() string {
	return "order_info"
}

// GetSaleByOrderId 获取根据订单id获取订单记录
func GetSaleByOrderId(OrderId int32) *OrderInfo {
	var oFind OrderInfo
	dao.DB.Model(&OrderInfo{}).Where("order_id=?", OrderId).First(&oFind)
	if oFind.OrderId == 0 {
		return nil
	}
	dFind := GetDrugByDrugId(oFind.DrugId)
	if dFind == nil {
		oFind.DrugName = "该药品信息已被删除"
	} else {
		oFind.DrugName = dFind.DrugName
	}
	return &oFind
}

// GetOrderListByUserId 获取根据用户id获取订单信息
func GetOrderListByUserId(userId int32) []*OrderInfo {
	var oListFind []*OrderInfo
	dao.DB.Model(&OrderInfo{}).Where("user_id=?", userId).Find(&oListFind)
	uFind := GetUserByUserId(userId)
	for _, s := range oListFind {
		if uFind == nil {
			s.UserName = "该用户信息已被删除"
		} else {
			s.UserName = uFind.UserName
		}
		dFind := GetDrugByDrugId(s.DrugId)
		if dFind == nil {
			s.DrugName = "该药品信息已被删除"
		} else {
			s.DrugName = dFind.DrugName
		}
	}
	return oListFind
}

// GetOrderListByDrugId 获取根据药品id获取订单信息
func GetOrderListByDrugId(drugId int32) []*OrderInfo {
	var oListFind []*OrderInfo
	dao.DB.Model(&OrderInfo{}).Where("drug_id=?", drugId).Find(&oListFind)
	for _, o := range oListFind {
		dFind := GetDrugByDrugId(o.DrugId)
		uFind := GetUserByUserId(o.UserId)
		if uFind == nil {
			o.UserName = "该用户信息已被删除"
		} else {
			o.UserName = uFind.UserName
		}
		if dFind == nil {
			o.DrugName = "该药品信息已被删除"
		} else {
			o.DrugName = dFind.DrugName
		}
	}
	return oListFind
}

// LikeGetOrderListByUserName 获取根据用户名称模糊查询订单列表
func LikeGetOrderListByUserName(userName string) []*OrderInfo {
	// 获取用户id
	uListFind := LikeGetUserListByUserName(userName)
	oListFind := []*OrderInfo{}
	for _, u := range uListFind {
		sList := GetOrderListByUserId(u.UserId)
		if len(sList) > 0 {
			oListFind = append(oListFind, sList...)
		}
	}
	return oListFind
}

// LikeGetOrderListByDrugName 获取根据药品名称模糊查询订单列表
func LikeGetOrderListByDrugName(drugName string) []*OrderInfo {
	// 获取药品id
	dListFind := LikeGetDrugListByDrugName(drugName)
	oListFind := []*OrderInfo{}
	for _, d := range dListFind {
		sList := GetOrderListByDrugId(d.DrugId)
		if len(sList) > 0 {
			oListFind = append(oListFind, sList...)
		}
	}
	return oListFind
}

// CreateOrder 创建订单记录
func CreateOrder(o *OrderInfo) (*OrderInfo, error) {
	if o == nil {
		return nil, errors.New("空指针错误")
	}
	err := dao.DB.Create(o).Error
	if err != nil {
		return nil, err
	}
	return o, nil
}

// GetOrderListByTime 获取根据日期查询订单列表
func GetOrderListByTime(startDate string, endDate string) []*OrderInfo {
	oListFind := []*OrderInfo{}
	if startDate != "" && endDate != "" {
		dao.DB.Where("create_time BETWEEN ? AND ?", startDate, endDate).Where("order_status=1").Find(&oListFind)
	} else {
		dao.DB.Find(&oListFind)
	}
	for _, s := range oListFind {
		uFind := GetUserByUserId(s.UserId)
		if uFind == nil {
			s.UserName = "该用户信息已被删除"
		} else {
			s.UserName = uFind.UserName
		}
		dFind := GetDrugByDrugId(s.DrugId)
		if dFind == nil {
			s.DrugName = "该药品信息已被删除"
		} else {
			s.DrugName = dFind.DrugName
		}
	}
	return oListFind
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(orderId int32, orderStatus int32) {
	dao.DB.Model(&OrderInfo{}).Where("order_id = ?", orderId).Update("order_status", orderStatus)
}

// UpdateOrderSupplyAmount 更新订单供应价格
func UpdateOrderSupplyAmount(orderId int32, supplyAmount decimal.Decimal) {
	dao.DB.Model(&OrderInfo{}).Where("order_id = ?", orderId).Update("supply_amount", supplyAmount)
}

// UpdateOrderStockInfo 更新订单库存信息
func UpdateOrderStockInfo(orderId int32, stockInfo string) {
	dao.DB.Model(&OrderInfo{}).Where("order_id = ?", orderId).Update("stock_info", stockInfo)
}
