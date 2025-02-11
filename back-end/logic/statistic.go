package logic

import (
	"drugims/model"
	"fmt"
	"sort"
	"time"
)

type PriceStatistics struct {
	Sale   float32 `json:"sale"`
	Profit float32 `json:"profit"`
	Supply float32 `json:"supply"`
}

type Statistics struct {
	TotalStatistics *PriceStatistics            `json:"total_statistics"`
	DailyStatistics map[string]*PriceStatistics `json:"daily_statistics"`
	DrugList        []*model.DrugInfo           `json:"drug_list"`
}

// GetStatisticByTime 获取统计信息
func GetStatisticByTime(startDate string, endDate string) (*Statistics, error) {
	statistics := &Statistics{
		TotalStatistics: &PriceStatistics{Sale: 0, Profit: 0, Supply: 0},
		DailyStatistics: make(map[string]*PriceStatistics),
		DrugList:        make([]*model.DrugInfo, 0),
	}

	// 获取指定日期内的销售信息和供应信息
	saleList := model.GetSaleListByTime(startDate, endDate)
	supplyList := model.GetSupplyListByTime(startDate, endDate)
	// 计算售出总价、进货总价以及总利润，并计算每种药的购买数量
	drugSaleMap := make(map[int32]int32)
	maxTime := time.Time{}
	minTime := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	for _, sale := range saleList {
		// 计算总的
		statistics.TotalStatistics.Sale += sale.SaleAmount
		statistics.TotalStatistics.Profit += sale.SaleAmount - sale.SupplyAmount

		// 计算每日的
		if sale.CreateTime.Before(minTime) {
			minTime = sale.CreateTime
		}
		if sale.CreateTime.After(maxTime) {
			maxTime = sale.CreateTime
		}
		date := fmt.Sprintf("%d-%02d-%02d", sale.CreateTime.Year(), sale.CreateTime.Month(), sale.CreateTime.Day())
		if _, ok := statistics.DailyStatistics[date]; !ok {
			statistics.DailyStatistics[date] = &PriceStatistics{Sale: 0, Profit: 0, Supply: 0}
		}
		statistics.DailyStatistics[date].Sale += sale.SaleAmount
		statistics.DailyStatistics[date].Profit += sale.SaleAmount - sale.SupplyAmount

		drugSaleMap[sale.DrugId] = sale.SaleQuantity
	}
	for _, supply := range supplyList {
		// 计算总的
		statistics.TotalStatistics.Supply += supply.SupplyPrice * float32(supply.SupplyQuantity)

		// 计算每日的
		if supply.CreateTime.Before(minTime) {
			minTime = supply.CreateTime
		}
		if supply.CreateTime.After(maxTime) {
			maxTime = supply.CreateTime
		}
		date := fmt.Sprintf("%d-%02d-%02d", supply.CreateTime.Year(), supply.CreateTime.Month(), supply.CreateTime.Day())
		if _, ok := statistics.DailyStatistics[date]; !ok {
			statistics.DailyStatistics[date] = &PriceStatistics{Sale: 0, Profit: 0, Supply: 0}
		}
		statistics.DailyStatistics[date].Supply += supply.SupplyPrice * float32(supply.SupplyQuantity)
	}

	// 遍历日期，将没数据的日期设为0
	s, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s = minTime
	}
	e, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		e = maxTime
	}
	for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if _, exists := statistics.DailyStatistics[dateStr]; !exists {
			statistics.DailyStatistics[dateStr] = &PriceStatistics{
				Sale:   0,
				Profit: 0,
				Supply: 0,
			}
		}
	}

	// 获取所有药品信息，按售出数量进行降序排序，前端根据出售数量和库存上下限阈值作比较给出需要变更阈值的提示
	statistics.DrugList = model.LikeGetDrugListByDrugName("")
	for _, drug := range statistics.DrugList {
		drug.SaleQuantity = drugSaleMap[drug.DrugId]
	}
	sort.Slice(statistics.DrugList, func(i, j int) bool {
		return statistics.DrugList[i].SaleQuantity > statistics.DrugList[j].SaleQuantity
	})

	return statistics, nil
}
