package sparkclient

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiVersion = "v3.1"
)

// ErrEmptyResponse is returned when the OpenAI API returns an empty response.
var ErrEmptyResponse = errors.New("empty response")

type APIVersion string

const (
	APIv1 APIVersion = "v1.1"
	APIv2 APIVersion = "v2.1"
	APIv3 APIVersion = "v3.1"
)

// Client is a client for the OpenAI API.
type Client struct {
	appId        string
	apiKey       string
	apiSecret    string
	Model        string
	baseURL      string
	organization string
	wsClient     Doer
	domain       string
	// required when APIVersion
	apiVersion      APIVersion
	embeddingsModel string
}

// Option is an option for the Spark client.
type Option func(*Client) error

// Doer performs a HTTP request.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New returns a new SparkAI client.
func New(domain, apiKey, apiSecret, appId string, baseURL string, organization string,
	apiVersion string, embeddingsModel string,
	opts ...Option,
) (*Client, error) {
	c := &Client{
		domain:          domain,
		apiSecret:       apiSecret,
		apiKey:          apiKey,
		appId:           appId,
		embeddingsModel: embeddingsModel,
		baseURL:         strings.TrimSuffix(baseURL, "/"),
		organization:    organization,
		apiVersion:      APIVersion(apiVersion),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Completion is a completion.
type Completion struct {
	Text string `json:"text"`
}

// CreateCompletion creates a completion.
func (c *Client) CreateCompletion(ctx context.Context, r *CompletionRequest) (messages.ChatMessage, error) {
	return c.createCompletion(ctx, r)

}

// EmbeddingRequest is a request to create an embedding.
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// CreateEmbedding creates embeddings.
//func (c *Client) CreateEmbedding(ctx context.Context, r *EmbeddingRequest) ([][]float32, error) {
//	if r.Model == "" {
//		r.Model = defaultEmbeddingModel
//	}
//
//	resp, err := c.createEmbedding(ctx, &embeddingPayload{
//		Model: r.Model,
//		Input: r.Input,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	if len(resp.Data) == 0 {
//		return nil, ErrEmptyResponse
//	}
//
//	embeddings := make([][]float32, 0)
//	for i := 0; i < len(resp.Data); i++ {
//		embeddings = append(embeddings, resp.Data[i].Embedding)
//	}
//
//	return embeddings, nil
//}

// CreateChat creates chat request.
func (c *Client) CreateChat(ctx context.Context, r *ChatRequest) (messages.ChatMessage, error) {

	resp, err := c.createChat(ctx, r, nil)
	if err != nil {
		return nil, err
	}
	switch resp.GetType() {
	case messages.ChatMessageTypeAI:
		fcMsg := resp.(messages.AIChatMessage)
		fmt.Println(fcMsg.FunctionCall.GetContent())
	default:
		fmt.Println(resp.GetContent())

	}

	return resp, nil
}

// CreateChat creates chat request.
func (c *Client) CreateChatWithCallBack(ctx context.Context, r *ChatRequest, stream_cb func(msg messages.ChatMessage) error) (messages.ChatMessage, error) {

	resp, err := c.createChat(ctx, r, stream_cb)
	if err != nil {
		return nil, err
	}
	if len(resp.GetContent()) == 0 {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

// 创建鉴权url  apikey 即 hmac username
func (c *Client) assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := c.HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func (c *Client) HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func (c *Client) buildURL(suffix string, model string) string {
	// spark ai implement:
	return fmt.Sprintf("%s%s", c.baseURL, suffix)
}
