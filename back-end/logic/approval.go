package logic

import (
	"drugims/model"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// ApprovalOperate 审批操作
func ApprovalOperate(approvalInfo *model.ApprovalInfo) error {
	switch approvalInfo.ApprovalType {
	case 0: // 退款审批
		orderInfo := &model.OrderInfo{}
		stockInfo := make([]*model.StockInfo, 0)
		if err := json.Unmarshal([]byte(approvalInfo.ApprovalContent), orderInfo); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(orderInfo.StockInfo), &stockInfo); err != nil {
			return err
		}
		u := model.QueryUserByUserId(orderInfo.UserId)
		// 删除冻结余额
		u.BlockBalance = u.BlockBalance.Sub(orderInfo.SaleAmount)
		if approvalInfo.ApprovalStatus == 1 { // 同意审批
			// 退回库存
			for _, stock := range stockInfo {
				s := model.GetStockByStockId(stock.StockId)
				if s == nil {
					model.CreateStock(stock)
				} else {
					s.RemainingQuantity += stock.RemainingQuantity
					model.UpdateStock(s)
				}
			}
			// 解冻余额
			u.Balance = u.Balance.Add(orderInfo.SaleAmount)
			// 修改订单状态
			model.UpdateOrderStatus(orderInfo.OrderId, 5)
		} else if approvalInfo.ApprovalStatus == 1 { // 拒绝审批
			// 修改订单状态
			model.UpdateOrderStatus(orderInfo.OrderId, 1)
		}
		model.UpdateUserInfo(u.UserId, u)
	case 1: // 供应审批
		supplyOrder := &model.SupplyOrder{}
		if err := json.Unmarshal([]byte(approvalInfo.ApprovalContent), supplyOrder); err != nil {
			return err
		}
		u := model.QueryUserByUserId(supplyOrder.UserId)
		// 删除冻结余额
		u.BlockBalance = u.BlockBalance.Sub(supplyOrder.SupplyPrice.Mul(decimal.NewFromInt(int64(supplyOrder.SupplyQuantity))))
		// 同意审批
		if approvalInfo.ApprovalStatus == 1 {
			// 存储库存信息
			if err := model.CreateStock(&model.StockInfo{
				DrugId:            supplyOrder.DrugId,
				BatchNumber:       supplyOrder.BatchNumber,
				ProductionDate:    supplyOrder.ProductionDate,
				SupplyPrice:       supplyOrder.SupplyPrice,
				RemainingQuantity: supplyOrder.SupplyQuantity,
			}); err != nil {
				return err
			}
			// 解冻资金转移
			u.Balance = u.Balance.Add(supplyOrder.SupplyPrice.Mul(decimal.NewFromInt(int64(supplyOrder.SupplyQuantity))))
			// 修改进货状态
			model.UpdateSupplyStatus(supplyOrder.SupplyId, 1)
		} else if approvalInfo.ApprovalStatus == 2 { // 拒绝审批
			// 修改进货状态
			model.UpdateSupplyStatus(supplyOrder.SupplyId, 2)
		}
		model.UpdateUserInfo(u.UserId, u)
	default:
		return fmt.Errorf("不存在的审批类型")
	}
	// 更新审批单
	model.UpdateApprovalInfo(approvalInfo)
	return nil
}
