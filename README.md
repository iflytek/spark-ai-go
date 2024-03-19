# 讯飞星火大模型接入库 (spark-ai-go)

本Go SDK库帮助用户更快体验讯飞星火大模型

## 项目地址

* Github: [https://github.com/iflytek/spark-ai-go](https://github.com/iflytek/spark-ai-python)
  欢迎点赞，star

## 前言

感谢开源的力量，希望讯飞开源越做越好，星火大模型效果越来越好！。


![!img](logo.png)

**本logo出自[星火大模型](https://xinghuo.xfyun.cn/)**

***感谢社区(LangchainGo项目以及SparkLLM部分committer)[项目正在开发中]***

## 近期规划新特性[待演进]

- [x] 极简的接入,快速调用讯飞星火大模型
- [ ] LLM类统一接口，快速切换业界大模型
- [x] Python版本[SDK](https://github.com/iflytek/spark-ai-python/)进行中


## 如何使用

### 示例代码

* 前置条件
  需要在 xfyun.cn 申请有权限的
    * app_id
    * api_key
    * api_secret

* URL/Domain配置请查看[doc](https://www.xfyun.cn/doc/spark/Web.html#_1-%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E)
* 运行测试脚本需要提前将 `.env.example` 拷贝为 `.env`并配置其中变量

内如如下:

```tpl
# spark 授权信息
SPARKAI_APP_ID=
SPARKAI_API_KEY=
SPARKAI_API_SECRET=
SPARKAI_DOMAIN=
SPARKAI_URL=wss://spark-api.xf-yun.com/v3.5/chat
```

### 一次性返回结果(非流式)


```golang
func TestLLMClientStream(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, err := spark.NewClient(spark.WithBaseURL(SPARK_API_BASE), spark.WithApiKey(SPARK_API_KEY), spark.WithApiSecret(SPARK_API_SECRET), spark.WithAppId(SPARK_APP_ID), spark.WithAPIDomain(SPARK_DOMAIN))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我润色并简化这句话: 国内科技大厂，创业公司都在演进大模型AI Agent且AI Agent框架仍然处于发展初期，采用开源化的AIAgent演进路线有助于快速构建影响力\n",
			},
		},
	}
    res, err := client.CreateChat(ctx, r)
    if err != nil {
        fmt.Print(err.Error())
        return
    }
    fmt.Print(res.GetContent())

}
```

### 流式返回结果


```golang
func TestLLMClientStream(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, err := spark.NewClient(spark.WithBaseURL(SPARK_API_BASE), spark.WithApiKey(SPARK_API_KEY), spark.WithApiSecret(SPARK_API_SECRET), spark.WithAppId(SPARK_APP_ID), spark.WithAPIDomain(SPARK_DOMAIN))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我润色并简化这句话: 国内科技大厂，创业公司都在演进大模型AI Agent且AI Agent框架仍然处于发展初期，采用开源化的AIAgent演进路线有助于快速构建影响力\n",
			},
		},
	}
	_, err = client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
		return
	}

}
```
其中可以在CreateWithCallBack传入一个func用于回调处理流式数据

### FunctionCall功能


```golang
func TestLLMFunctionCN(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: " 帮我生成一份2023ATP报表放到 /workspace目录下 ",
			},
		},
		Functions: []messages.FunctionDefinition{
			{
				Name:        "generate_biz_report",
				Description: "生成运营分析报表工具",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"biz": map[string]any{
							"type":        "string",
							"description": "业务类型,取值范围为 ase, aipaas, atp",
						},
						"output_dir": map[string]any{
							"type":        "string",
							"description": "生成存放位置",
						},
						"type": map[string]any{
							"type":        "string",
							"description": "报表类型，取值范围是 txt, excel, markdown之一,默认为excel",
						},
						"start": map[string]any{
							"type":        "string",
							"description": "起始时间年份",
						},
						"end": map[string]any{
							"type":        "string",
							"description": "结束时间年份",
						},
					},
					"required": []string{
						"output_dir", "type", "start", "end",
					},
				},
			},
			{
				Name:        "translate_transcript",
				Description: "using translate_text function to translate the script",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"source_language": map[string]any{
							"type":        "string",
							"description": "source language",
						},
						"target_language": map[string]any{
							"type":        "string",
							"description": "target language",
						},
					},
					"required": []string{
						"source_language",
						"source_language",
					},
				},
			},
		},
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		return
	}

}
```

## 欢迎贡献

扫码加入交流群

![img](weichat.jpg)

## 已知问题

* 项目目前开发阶段，有一些冗余代码，人力有限，部分思想借鉴开源实现

## URL

* wss://spark-api.xf-yun.com/v3.5/chat

## 致谢

* [LangchainGo](https://github.com/tmc/langchaingo)
