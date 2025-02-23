package logic

import (
	"drugims/model"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// SaleInfoFlow 销售信息流
type SaleInfoFlow struct {
	DrugInfo  *model.DrugInfo
	OrderInfo *model.OrderInfo
}

// 购买药品
func BuyDrug(drugInfo *model.DrugInfo) (*SaleInfoFlow, error) {
	return NewOrderInfoFlow(drugInfo, nil).buyDrug()
}

// 创建订单
func CreateOrder(drugInfo *model.DrugInfo) (*SaleInfoFlow, error) {
	return NewOrderInfoFlow(drugInfo, nil).createOrder()
}

// 确认订单
func ConfirmOrder(orderInfo *model.OrderInfo) (*SaleInfoFlow, error) {
	return NewOrderInfoFlow(nil, orderInfo).confirmOrder()
}

// 撤销订单
func RevokeOrder(orderInfo *model.OrderInfo) (*SaleInfoFlow, error) {
	return NewOrderInfoFlow(nil, orderInfo).revokeOrder()
}

// 退款
func RefundOrder(orderInfo *model.OrderInfo) (*SaleInfoFlow, error) {
	return NewOrderInfoFlow(nil, orderInfo).refundOrder()
}

func NewOrderInfoFlow(drugInfo *model.DrugInfo, orderInfo *model.OrderInfo) *SaleInfoFlow {
	return &SaleInfoFlow{DrugInfo: drugInfo, OrderInfo: orderInfo}
}

// 购买并消耗库存
func (s *SaleInfoFlow) buyDrug() (*SaleInfoFlow, error) {
	u := model.QueryUserByUserId(s.DrugInfo.UserId)
	saleTotalPrice := s.DrugInfo.SalePrice.Mul(decimal.NewFromInt(int64(s.DrugInfo.SaleQuantity)))
	switch s.DrugInfo.PaymentType {
	case 1:
		// 微信支付检查...
		u.BlockBalance = u.BlockBalance.Add(saleTotalPrice)
	case 2:
		// 支付宝支付检查...
		u.BlockBalance = u.BlockBalance.Add(saleTotalPrice)
	case 3:
		// 余额支付
		if u.Balance.LessThan(saleTotalPrice) {
			return nil, fmt.Errorf("您的余额不足, 请充值")
		}
		u.Balance = u.Balance.Sub(saleTotalPrice)
		u.BlockBalance = u.BlockBalance.Add(saleTotalPrice)
	}
	model.UpdateUserInfo(u.UserId, u)

	// 取出库存
	stockList := model.GetStockListByDrugId(s.DrugInfo.DrugId)

	// 根据生产日期从小到大排序
	sort.Slice(stockList, func(i, j int) bool {
		dateI, err1 := time.Parse("2006-01-02", stockList[i].ProductionDate)
		dateJ, err2 := time.Parse("2006-01-02", stockList[j].ProductionDate)
		if err1 != nil || err2 != nil {
			return stockList[i].StockId > stockList[j].StockId
		}
		return dateI.Before(dateJ)
	})

	// 扣之前再检查一遍
	totalStock := int32(0)
	for _, stock := range stockList {
		totalStock += stock.RemainingQuantity
	}
	// 库存不足则直接将订单状态设置为已关闭
	if totalStock < s.DrugInfo.SaleQuantity {
		model.UpdateOrderStatus(s.DrugInfo.OrderId, 2)
		return nil, fmt.Errorf("库存不足")
	}

	// 从旧到新，一个个更新库存数量
	buyNum := s.DrugInfo.SaleQuantity // 购买数量
	supplyTotalPrice := decimal.Zero  // 采购总价
	stockInfo := []*model.StockInfo{}
	for _, stock := range stockList {
		stockTemp := stock
		if stock.RemainingQuantity > buyNum {
			stock.RemainingQuantity -= buyNum
			supplyTotalPrice = supplyTotalPrice.Add(stock.SupplyPrice.Mul(decimal.NewFromInt(int64(buyNum))))
			// 更新改库存的数量
			if err := model.UpdateStock(stock); err != nil {
				return nil, err
			}
			// 增加库存信息
			stockTemp.RemainingQuantity = buyNum
			stockInfo = append(stockInfo, stockTemp)
			break
		} else {
			buyNum -= stock.RemainingQuantity
			supplyTotalPrice = supplyTotalPrice.Add(stock.SupplyPrice.Mul(decimal.NewFromInt(int64(stock.RemainingQuantity))))
			// 直接删掉改库存数据
			if err := model.DeleteStock(stock.StockId); err != nil {
				return nil, err
			}
			// 增加库存信息
			stockInfo = append(stockInfo, stockTemp)
		}
	}

	// 变更订单状态为待确认，添加供应总价以及库存
	model.UpdateOrderStatus(s.DrugInfo.OrderId, 3)
	model.UpdateOrderSupplyAmount(s.DrugInfo.OrderId, supplyTotalPrice)
	stockInfoJson, err := json.Marshal(stockInfo)
	if err != nil {
		return nil, err
	}
	model.UpdateOrderStockInfo(s.DrugInfo.OrderId, string(stockInfoJson))
	return s, nil
}

// 创建订单
func (s *SaleInfoFlow) createOrder() (*SaleInfoFlow, error) {
	// 生成订单记录
	OrderInfo := &model.OrderInfo{
		DrugId:       s.DrugInfo.DrugId,
		UserId:       s.DrugInfo.UserId,
		SaleQuantity: s.DrugInfo.SaleQuantity,
		SaleAmount:   s.DrugInfo.SalePrice.Mul(decimal.NewFromInt(int64(s.DrugInfo.SaleQuantity))),
		OrderStatus:  0,
	}
	o, err := model.CreateOrder(OrderInfo)
	if err != nil {
		return nil, err
	}
	s.OrderInfo = o
	return s, nil
}

// 确认订单
func (s *SaleInfoFlow) confirmOrder() (*SaleInfoFlow, error) {
	// 先变更订单状态
	model.UpdateOrderStatus(s.OrderInfo.OrderId, 1)
	// 变更成功后再扣除用户冻结金额
	u := model.QueryUserByUserId(s.OrderInfo.UserId)
	u.BlockBalance = u.BlockBalance.Sub(s.OrderInfo.SaleAmount)
	model.UpdateUserInfo(u.UserId, u)
	return s, nil
}

// 撤销订单
func (s *SaleInfoFlow) revokeOrder() (*SaleInfoFlow, error) {
	model.UpdateOrderStatus(s.OrderInfo.OrderId, 2)
	return s, nil
}

// 订单退款
func (s *SaleInfoFlow) refundOrder() (*SaleInfoFlow, error) {
	// 变更订单状态为审批中
	model.UpdateOrderStatus(s.OrderInfo.OrderId, 4)
	// 添加冻结余额
	u := model.QueryUserByUserId(s.OrderInfo.UserId)
	u.BlockBalance = u.BlockBalance.Add(s.OrderInfo.SaleAmount)
	model.UpdateUserInfo(u.UserId, u)
	// 提交退款审批
	orderInfoJson, err := json.Marshal(s.OrderInfo)
	if err != nil {
		return nil, err
	}
	approvalInfo := &model.ApprovalInfo{
		UserId:          s.OrderInfo.UserId,
		ApprovalType:    0,
		ApprovalContent: string(orderInfoJson),
		Reason:          s.OrderInfo.Reason,
		ApprovalStatus:  0,
	}
	model.CreateApproval(approvalInfo)
	return s, nil
}
