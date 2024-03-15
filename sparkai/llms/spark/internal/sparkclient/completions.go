package sparkclient

import (
	"context"
	"github.com/iflytek/spark-ai-go/sparkai/llms"
)

// CompletionRequest is a request to complete a completion.
type CompletionRequest struct {
	Prompt      string                    `json:"prompt"`
	Temperature float64                   `json:"temperature,omitempty"`
	MaxTokens   int64                     `json:"max_tokens,omitempty"`
	N           int                       `json:"n,omitempty"`
	TopK        int64                     `json:"top_k,omitempty"`
	Functions   []llms.FunctionDefinition `json:"functions"`
}

type CompletionResponse struct {
	ID      string  `json:"id,omitempty"`
	Created float64 `json:"created,omitempty"`
	Choices []struct {
		FinishReason string      `json:"finish_reason,omitempty"`
		Index        float64     `json:"index,omitempty"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		Text         string      `json:"text,omitempty"`
	} `json:"choices,omitempty"`
	Model  string `json:"model,omitempty"`
	Object string `json:"object,omitempty"`
	Usage  struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
	FunctionCall *llms.FunctionCall `json:"function_call"`
}

type errorMessage struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// 生成参数
func (c *Client) constructSparkReq(appid string, req *ChatRequest) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := req.Messages
	if req.Domain == nil {
		req.Domain = &defaultDomain
	}
	if req.Temperature == nil {
		req.Temperature = &defaultTemperature
	}
	if req.TopK == nil {
		req.TopK = &defaultTopK
	}
	if req.MaxTokens == nil {
		req.MaxTokens = &defaultMaxTokens
	}
	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      req.Domain,      // 根据实际情况修改返回的数据结构和字段名
				"temperature": req.Temperature, // 根据实际情况修改返回的数据结构和字段名
				"top_k":       req.TopK,        // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  req.MaxTokens,   // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",       // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
			"functions": map[string][]llms.FunctionDefinition{
				"text": req.Functions,
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// nolint:lll
func (c *Client) createCompletion(ctx context.Context, payload *CompletionRequest) (llms.ChatMessage, error) {
	return c.createChat(ctx, &ChatRequest{
		Messages: []llms.ChatMessage{
			&(llms.GenericChatMessage{Role: "user", Content: payload.Prompt}),
		},
		Temperature: &payload.Temperature,
		TopK:        &payload.TopK,
		MaxTokens:   &payload.MaxTokens,
		Functions:   payload.Functions,
	}, nil)
}
