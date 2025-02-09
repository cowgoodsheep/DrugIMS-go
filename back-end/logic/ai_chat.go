package logic

import (
	"bytes"
	"drugims/config"
	"drugims/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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
			"售价:" + strconv.Itoa(int(drugInfo.SalePrice)) + "\n"
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
