package test

import (
	"context"
	"fmt"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/client/sparkclient"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")

}

const (
	AppIdEnvVarName        = "SPARKAI_APP_ID"     //nolint:gosec
	ApiKeyEnvVarName       = "SPARKAI_API_KEY"    //nolint:gosec
	ApiSecretEnvVarName    = "SPARKAI_API_SECRET" //nolint:gosec
	SparkDomainEnvVarName  = "SPARKAI_DOMAIN"
	sparkVersionEnvVarName = "SPARKAI_API_VERSION" //nolint:gosec
	BaseURLEnvVarName      = "SPARKAI_URL"         //nolint:gosec
	organizationEnvVarName = "SPARK_ORGANIZATION"  //nolint:gosec
)

func TestSpark(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	llm, err := spark.New(spark.WithAPIDomain(SPARK_DOMAIN),
		spark.WithApiKey(SPARK_API_KEY),
		spark.WithApiSecret(SPARK_API_SECRET),
		spark.WithAppId(SPARK_APP_ID),
		spark.WithBaseURL(SPARK_API_BASE))

	ctx := context.Background()
	g, err := llm.Generate(ctx, []string{"帮我润色并简化这句话: 国内科技大厂，创业公司都在演进大模型AI Agent且AI Agent框架仍然处于发展初期，采用开源化的AIAgent演进路线有助于快速构建影响力\n", "帮我赞扬一下库里的表现"})

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	for _, c := range g {
		fmt.Println(c)
	}

}

func TestLLMClientStream(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = ""
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

//func TestLLMFunctionCall(t *testing.T) {
//	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
//	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
//	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
//	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
//	SPARK_DOMAIN := "10245"
//	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
//	_, client, _ := newClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
//	ctx := context.Background()
//	r := &sparkclient.ChatRequest{
//		Domain: &SPARK_DOMAIN,
//		Messages: []messages.ChatMessage{
//			&messages.GenericChatMessage{
//				Role:    "user",
//				Content: " For the case: For the video located in /usr/local/3.mp4, recognize the speech and transfer it into a script file, please choose a function to complete it",
//			},
//		},
//		Functions: []messages.FunctionDefinition{
//			{
//				Name:        "recognize_transcript_from_video",
//				Description: "recognize the speech from video and transfer into a txt file",
//				Parameters: map[string]any{
//					"type": "object",
//					"properties": map[string]any{
//						"audio_filepath": map[string]any{
//							"type":        "string",
//							"description": "path of the vedio file",
//						},
//					},
//					"required": []string{
//						"audio_filepath",
//					},
//				},
//			},
//			{
//				Name:        "translate_transcript",
//				Description: "using translate_text function to translate the script",
//				Parameters: map[string]any{
//					"type": "object",
//					"properties": map[string]any{
//						"source_language": map[string]any{
//							"type":        "string",
//							"description": "source language",
//						},
//						"target_language": map[string]any{
//							"type":        "string",
//							"description": "target language",
//						},
//					},
//					"required": []string{
//						"source_language",
//						"source_language",
//					},
//				},
//			},
//		},
//	}
//	//_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
//	//	fmt.Print(msg.GetContent())
//	//	return nil
//	//})
//	rsp, err := client.CreateChat(ctx, r)
//	if err != nil {
//		fmt.Print(rsp.GetType(), rsp.GetContent())
//		return
//	}
//
//}
//
//func TestLLMFunctionCallWithCallBack(t *testing.T) {
//	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
//	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
//	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
//	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
//	SPARK_DOMAIN := "10245"
//	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
//	_, client, _ := newClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
//	ctx := context.Background()
//	r := &sparkclient.ChatRequest{
//		Domain: &SPARK_DOMAIN,
//		Messages: []messages.ChatMessage{
//			&messages.GenericChatMessage{
//				Role:    "user",
//				Content: " For the case: For the video located in /usr/local/3.mp4, recognize the speech and transfer it into a script file, please choose a function to complete it",
//			},
//		},
//		Functions: []messages.FunctionDefinition{
//			{
//				Name:        "recognize_transcript_from_video",
//				Description: "recognize the speech from video and transfer into a txt file",
//				Parameters: map[string]any{
//					"type": "object",
//					"properties": map[string]any{
//						"audio_filepath": map[string]any{
//							"type":        "string",
//							"description": "path of the vedio file",
//						},
//					},
//					"required": []string{
//						"audio_filepath",
//					},
//				},
//			},
//			{
//				Name:        "translate_transcript",
//				Description: "using translate_text function to translate the script",
//				Parameters: map[string]any{
//					"type": "object",
//					"properties": map[string]any{
//						"source_language": map[string]any{
//							"type":        "string",
//							"description": "source language",
//						},
//						"target_language": map[string]any{
//							"type":        "string",
//							"description": "target language",
//						},
//					},
//					"required": []string{
//						"source_language",
//						"source_language",
//					},
//				},
//			},
//		},
//	}
//	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
//		fmt.Print(msg.GetContent())
//		return nil
//	})
//	if err != nil {
//		return
//	}
//
//}
