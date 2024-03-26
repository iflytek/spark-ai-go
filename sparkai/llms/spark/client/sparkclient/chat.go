package sparkclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/iflytek/spark-ai-go/log"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	defaultDomain              = "general"
	defaultTemperature float64 = float64(0.8)
	defaultTopK                = int64(6)
	defaultMaxTokens           = int64(2048)
)

var ErrContentExclusive = errors.New("only one of Content / MultiContent allowed in message")

// ChatRequest is a request to complete a chat completion..
type ChatRequest struct {
	Domain      *string                       `json:"domain"`
	Messages    []messages.ChatMessage        `json:"messages"`
	Temperature *float64                      `json:"temperature,omitempty"`
	TopK        *int64                        `json:"top_p,omitempty"`
	MaxTokens   *int64                        `json:"max_tokens,omitempty"`
	Audit       *string                       `json:"audit,omitempty"`
	Functions   []messages.FunctionDefinition `json:"functions,omitempty"`

	//// Function definitions to include in the request.
	//// FunctionCallBehavior is the behavior to use when calling functions.
	////
	//// If a specific function should be invoked, use the format:
	//// `{"name": "my_function"}`
	//FunctionCallBehavior FunctionCallBehavior `json:"function_call,omitempty"`
	//
	//// StreamingFunc is a function to be called for each chunk of a streaming response.
	//// Return an error to stop streaming early.
	//StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct { //nolint:musttag
	// The role of the author of this message. One of system, user, or assistant.
	Role string
	// The content of the message.
	Content string

	//MultiContent []llms.ContentPart
	//
	//// The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,
	//// with a maximum length of 64 characters.
	//Name string

	// FunctionCall represents a function call to be made in the message.
	FunctionCall *messages.FunctionCall
}

func (m ChatMessage) GetType() messages.ChatMessageType {
	return messages.ChatMessageType(m.Role)
}

func (m ChatMessage) GetContent() string {
	return m.Content
}
func (m ChatMessage) MarshalJSON() ([]byte, error) {
	msg := struct {
		Role         string                 `json:"role"`
		Content      string                 `json:"content"`
		FunctionCall *messages.FunctionCall `json:"function_call"`
	}{
		Role:         m.Role,
		Content:      m.Content,
		FunctionCall: m.FunctionCall,
	}

	return json.Marshal(msg)
}

func (m *ChatMessage) UnmarshalJSON(data []byte) error {
	msg := struct {
		Role         string                 `json:"role"`
		Content      string                 `json:"content"`
		FunctionCall *messages.FunctionCall `json:"function_call"`
	}{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}
	*m = ChatMessage(msg)
	return nil
}

// ChatChoice is a choice in a chat response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatUsage is the usage of a chat completion request.
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatResponse is a response to a chat request.
type ChatResponse struct {
	//Choices []*ChatChoice `json:"choices,omitempty"`
	Role         string                 `json:"role"`
	Content      string                 `json:"content,omitempty"`
	FunctionCall *messages.FunctionCall `json:"function_call"`
	Usage        struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

func (c *ChatResponse) GetType() messages.ChatMessageType {
	if c.FunctionCall != nil {
		return messages.ChatMessageTypeFunction
	}
	return messages.ChatMessageType(c.Role)

}

func (c *ChatResponse) GetContent() string {
	return c.Content
}

func (c *ChatResponse) UpdateContent(msg string) {
	c.Content = msg
}

// StreamedChatResponsePayload is a chunk from the stream.
type StreamedChatResponsePayload struct {
	ID      string  `json:"id,omitempty"`
	Created float64 `json:"created,omitempty"`
	Model   string  `json:"model,omitempty"`
	Object  string  `json:"object,omitempty"`
	Choices []struct {
		Index float64 `json:"index,omitempty"`
		Delta struct {
			Role         string        `json:"role,omitempty"`
			Content      string        `json:"content,omitempty"`
			FunctionCall *FunctionCall `json:"function_call,omitempty"`
		} `json:"delta,omitempty"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
}

// FunctionDefinition is a definition of a function that can be called by the model.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of the function.
	Description string `json:"description"`
	// Parameters is a list of parameters for the function.
	Parameters any `json:"parameters"`
}

// FunctionCallBehavior is the behavior to use when calling functions.
type FunctionCallBehavior string

const (
	// FunctionCallBehaviorUnspecified is the empty string.
	FunctionCallBehaviorUnspecified FunctionCallBehavior = ""
	// FunctionCallBehaviorNone will not call any functions.
	FunctionCallBehaviorNone FunctionCallBehavior = "none"
	// FunctionCallBehaviorAuto will call functions automatically.
	FunctionCallBehaviorAuto FunctionCallBehavior = "auto"
)

// FunctionCall is a call to a function.
type FunctionCall struct {
	// Name is the name of the function to call.
	Name string `json:"name"`
	// Arguments is the set of arguments to pass to the function.
	Arguments string `json:"arguments"`
}

func (c *Client) createChat(ctx context.Context, payload *ChatRequest, cb func(msg messages.ChatMessage) error) (messages.ChatMessage, error) {

	// Build request payload

	// Build request
	if c.baseURL == "" {
		return nil, errors.New("No API Url set")
	}
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	ua_str := ""
	user_agent := ctx.Value("user_agent")
	if user_agent != nil {
		ua_str = user_agent.(string)
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(c.assembleAuthUrl1(c.baseURL, c.apiKey, c.apiSecret), map[string][]string{"User-Agent": []string{fmt.Sprintf("SparkAISdk/golang %s", ua_str)}})
	if err != nil {
		return nil, errors.New(err.Error())

	} else if resp.StatusCode != 101 {
		panic(readResp(resp) + err.Error())
	}

	go func() {

		data := c.constructSparkReq(c.appId, payload)
		conn.WriteJSON(data)

	}()
	var answer = ""
	var code int
	// Parse response
	var response messages.ChatMessage

	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}

		var sparkResp = messages.SparkResponse{}
		err1 := json.Unmarshal(msg, &sparkResp)
		if err1 != nil {
			return nil, errors.New(err1.Error())
		}
		//fmt.Println(string(msg))
		//解析数据
		//header := data["header"].(map[string]interface{})
		code = sparkResp.Header.Code

		if code != 0 {
			//fmt.Println(data["payload"])
			return nil, errors.New(fmt.Sprintf("code != 0 %f", code))

		}
		payload := &sparkResp.Payload
		//choices := payload["choices"].(map[string]interface{})
		choices := payload.Choices
		status := choices.Status
		text := choices.Text
		content := text[0].Content
		role := text[0].Role
		fc := text[0].FunctionCall
		if fc != nil {
			response = messages.AIChatMessage{
				Content:      fc.GetContent(),
				FunctionCall: fc,
			}
		} else {
			response = messages.GenericChatMessage{
				Content: content,
				Role:    role,
			}
		}
		// 处理 cb todo cb 规范化
		if cb != nil {

			err = cb(response)
			if err != nil {
				fmt.Println("callback error ")
			}
		}

		if status != 2 {
			answer += content
		} else {
			answer += content
			usage := payload.Usage
			//totalTokens = temp["total_tokens"].(float64)
			//fmt.Println("total_tokens:", totalTokens)
			if fc == nil {
				response = &ChatResponse{
					Role:         role,
					Content:      answer,
					FunctionCall: fc,
					Usage: struct {
						CompletionTokens float64 `json:"completion_tokens,omitempty"`
						PromptTokens     float64 `json:"prompt_tokens,omitempty"`
						TotalTokens      float64 `json:"total_tokens,omitempty"`
					}{
						CompletionTokens: usage.CompletionTokens,
						PromptTokens:     usage.PromptTokens,
						TotalTokens:      usage.TotalTokens,
					},
				}
			}
			log.GetLogger().Info("Sid: ", sparkResp.Header.Sid)

			conn.Close()
			break
		}

	}
	if code != 0 {
		msg := fmt.Sprintf("API returned unexpected status code: %f", code)

		// No need to check the error here: if it fails, we'll just return the
		// status code.

		return nil, fmt.Errorf("%s", msg) // nolint:goerr113
	}

	return response, nil
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
