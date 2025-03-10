package logic

import (
	"bytes"
	"drugims/config"
	"drugims/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type PriceStatistics struct {
	Sale   decimal.Decimal `json:"sale"`
	Profit decimal.Decimal `json:"profit"`
	Supply decimal.Decimal `json:"supply"`
}

type Statistics struct {
	TotalStatistics *PriceStatistics            `json:"total_statistics"`
	DailyStatistics map[string]*PriceStatistics `json:"daily_statistics"`
	DrugList        []*model.DrugInfo           `json:"drug_list"`
}

// GetStatisticByTime 获取统计信息
func GetStatisticByTime(startDate string, endDate string) (*Statistics, error) {
	statistics := &Statistics{
		TotalStatistics: &PriceStatistics{Sale: decimal.Zero, Profit: decimal.Zero, Supply: decimal.Zero},
		DailyStatistics: make(map[string]*PriceStatistics),
		DrugList:        make([]*model.DrugInfo, 0),
	}

	// 获取指定日期内的订单信息和供应信息
	orderList := model.GetOrderListByTime(startDate, endDate)
	supplyList := model.GetSupplyListByTime(startDate, endDate)
	// 计算售出总价、进货总价以及总利润，并计算每种药的购买数量
	drugSaleMap := make(map[int32]int32)
	maxTime := time.Time{}
	minTime := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	for _, order := range orderList {
		if order.OrderStatus == 1 { // 只计算已完成的订单
			// 计算总的
			statistics.TotalStatistics.Sale = statistics.TotalStatistics.Sale.Add(order.SaleAmount)
			statistics.TotalStatistics.Profit = statistics.TotalStatistics.Profit.Add(order.SaleAmount.Sub(order.SupplyAmount))

			// 计算每日的
			if order.CreateTime.Before(minTime) {
				minTime = order.CreateTime
			}
			if order.CreateTime.After(maxTime) {
				maxTime = order.CreateTime
			}
			date := fmt.Sprintf("%d-%02d-%02d", order.CreateTime.Year(), order.CreateTime.Month(), order.CreateTime.Day())
			if _, ok := statistics.DailyStatistics[date]; !ok {
				statistics.DailyStatistics[date] = &PriceStatistics{Sale: decimal.Zero, Profit: decimal.Zero, Supply: decimal.Zero}
			}
			statistics.DailyStatistics[date].Sale = statistics.DailyStatistics[date].Sale.Add(order.SaleAmount)
			statistics.DailyStatistics[date].Profit = statistics.DailyStatistics[date].Profit.Add(order.SaleAmount.Sub(order.SupplyAmount))

			drugSaleMap[order.DrugId] = order.SaleQuantity
		}
	}
	for _, supply := range supplyList {
		// 计算总的
		statistics.TotalStatistics.Supply = statistics.TotalStatistics.Supply.Add(supply.SupplyPrice.Mul(decimal.NewFromInt(int64(supply.SupplyQuantity))))

		// 计算每日的
		if supply.CreateTime.Before(minTime) {
			minTime = supply.CreateTime
		}
		if supply.CreateTime.After(maxTime) {
			maxTime = supply.CreateTime
		}
		date := fmt.Sprintf("%d-%02d-%02d", supply.CreateTime.Year(), supply.CreateTime.Month(), supply.CreateTime.Day())
		if _, ok := statistics.DailyStatistics[date]; !ok {
			statistics.DailyStatistics[date] = &PriceStatistics{Sale: decimal.Zero, Profit: decimal.Zero, Supply: decimal.Zero}
		}
		statistics.DailyStatistics[date].Supply = statistics.DailyStatistics[date].Supply.Add(supply.SupplyPrice.Mul(decimal.NewFromInt(int64(supply.SupplyQuantity))))
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
				Sale:   decimal.Zero,
				Profit: decimal.Zero,
				Supply: decimal.Zero,
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

// ----------------------------------------------------------------
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens            int         `json:"prompt_tokens"`
		CompletionTokens        int         `json:"completion_tokens"`
		TotalTokens             int         `json:"total_tokens"`
		CompletionTokensDetails interface{} `json:"completion_tokens_details"`
		PromptTokensDetails     interface{} `json:"prompt_tokens_details"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

// GetAiChatResponse 获取ai回答
func GetAiChatResponse(questionList []Message) (string, error) {
	// 构造请求数据
	reqData := ChatRequest{
		Model:    config.Conf.AiChat.Model,
		Messages: []Message{},
	}
	drugInfoList := model.LikeGetDrugListByDrugName("")
	drugInfoMsg := "你是一个专业的医生,药店的药物信息如下：\n"
	for _, drugInfo := range drugInfoList {
		drugInfoMsg += "名称:" + drugInfo.DrugName + "," +
			"说明:" + drugInfo.DrugDescription + "," +
			"厂家:" + drugInfo.Manufacturer + "," +
			"规格:" + drugInfo.Specification + "," +
			"库存下限:" + strconv.Itoa(int(drugInfo.StockLowerLimit)) + "," +
			"库存上限:" + strconv.Itoa(int(drugInfo.StockUpperLimit)) + "," +
			"库存剩余:" + strconv.Itoa(int(drugInfo.StockRemain)) + "," +
			"售价:" + drugInfo.SalePrice.String() + "\n"
	}
	drugInfoMsg += "根据用户的病情做出诊断并推荐药物,要根据药店已有药品的实际情况回答"
	reqData.Messages = append(reqData.Messages, Message{Role: "system", Content: drugInfoMsg})
	reqData.Messages = append(reqData.Messages, questionList...)

	fmt.Println(reqData)

	// 将请求数据序列化为 JSON
	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", config.Conf.AiChat.URL, bytes.NewBuffer(reqDataJson))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.Conf.AiChat.Key)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应内容
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

// ----------------------------------------------------------------
type WarningInfo struct {
	WarningTitle       string `json:"warning_title"`       // 警告标题
	WarningDescription string `json:"warning_description"` // 警告内容
}

func RiskManage() []*WarningInfo {
	WarningInfoList := []*WarningInfo{}
	// 1. 低价检查
	drugInfoList := model.LikeGetDrugListByDrugName("")
	for _, drug := range drugInfoList {
		if drug.SalePrice.LessThan(decimal.NewFromInt(5)) {
			warningInfo := &WarningInfo{
				WarningTitle:       "低价警告",
				WarningDescription: fmt.Sprintf("药品ID: %d | 药品名称: %s | 目前售价: %s CNY", drug.DrugId, drug.DrugName, drug.SalePrice),
			}
			WarningInfoList = append(WarningInfoList, warningInfo)
		}
	}
	// 2. 反洗钱检查 todo

	return WarningInfoList
}
