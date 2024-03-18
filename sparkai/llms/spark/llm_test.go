package spark

import (
	"context"
	"fmt"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/client/sparkclient"
	"github.com/iflytek/spark-ai-go/sparkai/memory/file_memory"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"os"
	"testing"
)

func TestLLMClient(t *testing.T) {

	SPARK_APP_ID := ""
	SPARK_API_KEY := ""
	SPARK_API_SECRET := ""
	SPARK_DOMAIN := "generalv3.5"
	SPARK_API_BASE := "wss://spark-api.xf-yun.com/v3.5/chat"

	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: "请你模拟詹姆斯的口吻写一篇转会消息",
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

func TestLLMClientStream(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
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
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
		return
	}

}

func TestLLMFunctionCall(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: " For the case: For the video located in /usr/local/3.mp4, recognize the speech and transfer it into a script file, please choose a function to complete it",
			},
		},
		Functions: []messages.FunctionDefinition{
			{
				Name:        "recognize_transcript_from_video",
				Description: "recognize the speech from video and transfer into a txt file",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"audio_filepath": map[string]any{
							"type":        "string",
							"description": "path of the vedio file",
						},
					},
					"required": []string{
						"audio_filepath",
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
	//_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
	//	fmt.Print(msg.GetContent())
	//	return nil
	//})
	rsp, err := client.CreateChat(ctx, r)
	if err != nil {
		fmt.Print(rsp.GetType(), rsp.GetContent())
		return
	}

}

func TestLLMFunctionCallWithCallBack(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: " For the case: For the video located in /usr/local/3.mp4, recognize the speech and transfer it into a script file, please choose a function to complete it",
			},
		},
		Functions: []messages.FunctionDefinition{
			{
				Name:        "recognize_transcript_from_video",
				Description: "recognize the speech from video and transfer into a txt file",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"audio_filepath": map[string]any{
							"type":        "string",
							"description": "path of the vedio file",
						},
					},
					"required": []string{
						"audio_filepath",
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

func TestLLMFunctionCN(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
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

func TestLLMFunctionNest(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
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
							"type":        "object",
							"description": "业务属性字段",
							"properties": map[string]any{
								"type": map[string]any{
									"type":        "string",
									"description": "报表类型，取值范围是 txt, excel, markdown之一,默认为excel",
								},
							},
						},
						"output_dir": map[string]any{
							"type":        "string",
							"description": "生成存放位置",
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

func TestLLMClientSystemMsgSingle(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: "当前你是一个辩论赛主持人角色，正在进行一场题目为大学生该不该谈恋爱的辩论",
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

var systems = []string{
	"You are a helpful AI assistant.\nSolve tasks using your coding and language skills.\nIn the following cases, suggest python code (in a python coding block) or shell script (in a sh coding block) for the user to execute.\n    1. When you need to collect info, use the code to output the info you need, for example, browse or search the web, download/read a file, print the content of a webpage or a file, get the current date/time, check the operating system. After sufficient info is printed and the task is ready to be solved based on your language skill, you can solve the task by yourself.\n    2. When you need to perform some task with code, use the code to perform the task and output the result. Finish the task smartly.\nSolve the task step by step if you need to. If a plan is not provided, explain your plan first. Be clear which step uses code, and which step uses your language skill.\nWhen using code, you must indicate the script type in the code block. The user cannot provide any other feedback or perform any other action beyond executing the code you suggest. The user can't modify your code. So do not suggest incomplete code which requires users to modify. Don't use a code block if it's not intended to be executed by the user.\nIf you want the user to save the code in a file before executing it, put # filename: <filename> inside the code block as the first line. Don't include multiple code blocks in one response. Do not ask users to copy and paste the result. Instead, use 'print' function for the output when relevant. Check the execution result returned by the user.\nIf the result indicates there is an error, fix the error and output the code again. Suggest the full code instead of partial code or code changes. If the error can't be fixed or if the task is not solved even after the code is executed successfully, analyze the problem, revisit your assumption, collect additional info you need, and think of a different approach to try.\nWhen you find an answer, verify the answer carefully. Include verifiable evidence in your response if possible.\nReply \"TERMINATE\" in the end when everything is done",
	"你是一个天气助手,请根据我的输入为我返回相应的天气助手调用参数",
	"你是一个任务规划助手,专门负责任务规划生成。请根据我的输入内容合理拆解返回个规划列表，规划列表用一个JSON数组返回，形如 [\"任务1内容\"，\"任务2内容\"]，确保数组可以被 python json.loads解析,如果无法生成请返回空数组 []",
	"你是一个任务规划助手,请根据我提供的工具集合生成一份调用工具的规划列表，形如 [\"工具1名称\"，\"工具2名称\"]，确保数组可以被 python json.loads解析,如果无法生成请返回空数组 []",
	"当前已经执行的工具:  get_time",
}
var users = []string{
	"帮我写一个贪吃蛇游戏",
	"帮我查询下去年今天合肥市的天气, 当前的工具有 get_weather(查询天气插件),get_time(查询当前时间)", // 如果仅仅是稳查询天梯，而不提供工具信息，大模型规划容易没有颗粒度，他不知道用那个什么颗粒度去拆分子任务
}
var assistants = []string{
	"[\"get_time\", \"get_weather\"]",
}

var weather = []messages.FunctionDefinition{
	{
		Name:        "get_weather",
		Description: "查询天气插件",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"province": map[string]any{
					"type":        "string",
					"description": "省份信息",
				},
				"city": map[string]any{
					"type":        "string",
					"description": "城市信息",
				},
				"time": map[string]any{
					"type":        "string",
					"description": "时间信息，格式必须为  2023-01-10 18:30:20",
				},
			},
			"required": []string{
				"province", "city", "time",
			},
		},
	}, {
		Name:        "get_time",
		Description: "查询当前时间",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"zone": map[string]any{
					"type":        "string",
					"description": "时区， 取值(utc, beijing) 默认beijing ,为北京时间",
				},
			},
			"required": []string{
				"zone",
			},
		},
	},
}

func TestLLMClientSystemMsg1Ground(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: systems[3],
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: users[1],
			},
		},
		//Functions: weather,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestLLMClientSystemMsg2Ground(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: "你是一个任务规划助手,请根据我提供的工具集合生成一份调用工具的规划列表，形如 [\"工具1名称\"，\"工具2名称\"]，确保数组可以被 python json.loads解析,如果无法生成请返回空",
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我查询下去年今天合肥市的天气, 当前的工具有 get_weather(查询天气插件),get_time(查询当前时间)",
			},
			messages.GenericChatMessage{
				Role:    "assistant",
				Content: "[\"get_time\", \"get_weather\"]",
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我查询下去年今天合肥市的天气, 当前的工具有 get_weather(查询天气插件),get_time(查询当前时间)",
			},
			messages.GenericChatMessage{
				Role:    "function",
				Content: "2024-01-23 13:00:33",
			},
			messages.GenericChatMessage{
				Role:    "system",
				Content: "当前已经执行过工具get_time, 它返回结果是: 2024-01-25 13:00:33",
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我查询下去年今天合肥市的天气, 当前的工具有 get_weather(查询天气插件),get_time(查询当前时间)",
			},
		},
		Functions: weather,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

var debate_systems = []string{

	"你是一个辩论赛主持人，请根据用户目标以及用户提供的方法集合生成一份调用方法的规划列表， 当前用户的方法有: gen_sub (根据题目生成辩论赛赛点), set_mode(设置辩论赛模式，支持1v1和3v3 2种模式 , set_max_round (设置最大轮数), select_speaker(选择一个辩论赛成员发表观点), summary (总结陈述阶段)。 返回的规划列表的元素只需要包含方名称,形如 [\"方法1名称\"，\"方法2名称\"]，确保数组可以被 python json.loads解析,如果无法生成请返回空数组 []。",
	"当前已经执行的工具:  get_time",
}
var debate_users = []string{
	"请以大学生该不该谈恋爱为题目,开始一场1v1的3轮辩论赛",
	// 如果仅仅是稳查询天梯，而不提供工具信息，大模型规划容易没有颗粒度，他不知道用那个什么颗粒度去拆分子任务
}
var debate_assistants = []string{
	"[\"get_time\", \"get_weather\"]",
}

var debate_methods = []messages.FunctionDefinition{
	{
		Name:        "get_sub",
		Description: "生成辩论观点",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"topic": map[string]any{
					"type":        "string",
					"description": "辩论赛话题",
				},
			},
			"required": []string{
				"topic",
			},
		},
	}, {
		Name:        "set_mode",
		Description: "设置辩论赛模式",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"mode": map[string]any{
					"type":        "string",
					"description": "辩论赛模式，支持 1v1和 3v3两种模式",
				},
			},
			"required": []string{
				"mode",
			},
		},
	},
	{
		Name:        "set_max_round",
		Description: "设置辩论赛最大轮数",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"round_num": map[string]any{
					"type":        "number",
					"description": "辩论赛最大轮数，最大不超过20",
				},
			},
			"required": []string{
				"mode",
			},
		},
	},
	{
		Name:        "select_speaker",
		Description: "选择下一个辩手",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"round_num": map[string]any{
					"type":        "number",
					"description": "辩论赛最大轮数，最大不超过20",
				},
			},
			"required": []string{
				"round_num",
			},
		},
	},
	{
		Name:        "summary",
		Description: "陈述",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		},
	},
}

func TestLLMDebate1Ground(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: debate_systems[0],
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: debate_users[0],
			},
		},
		Functions: debate_methods,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestLLMDebate2Ground(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: debate_systems[0],
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: debate_users[0],
			},
			messages.GenericChatMessage{
				Role:    "assistant",
				Content: "[\"gen_sub\", \"set_mode\", \"set_max_round\", \"select_speaker\", \"summary\"]",
			},
			messages.GenericChatMessage{
				Role:    "function",
				Content: "gen_sub生成论题结果为: 正方: 大学生谈恋爱好， 反方: 大学生谈恋爱不好",
			},
			messages.GenericChatMessage{
				Role:    "system",
				Content: "你是个辩论赛主持人，当前辩论赛议程为:\n [\"gen_sub\", \"set_mode\", \"set_max_round\", \"select_speaker\", \"summary\"],\n已经进行到 [\"gen_sub\"],请根据以上议程和已经执行过的议程，生成下一步需要调用的议程方法和参数",
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: "请以大学生该不该谈恋爱为题目,开始一场1v1的3轮辩论赛",
			},
		},
		Functions: debate_methods,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

var train_methods = []messages.FunctionDefinition{
	{
		Name:        "order_train",
		Description: "预订火车函数",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"start_time": map[string]any{
					"type":        "string",
					"description": "火车发车时间",
				},
				"arrived_time": map[string]any{
					"type":        "string",
					"description": "火车到达时间",
				},
				"start_city": map[string]any{
					"type":        "string",
					"description": "火车发车地点",
				},
				"dest_city": map[string]any{
					"type":        "string",
					"description": "火车目的地点",
				},
				"stragedy": map[string]any{
					"type":        "string",
					"description": "订票策略, 分为价格优先，时间最短优先",
				},
			},
			"required": []string{
				"start_time", "start_city", "dest_city", "stragedy", "arrived_time",
			},
		},
	},
}

var train_systems = []string{
	"你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE.",
}

func TestTrainOrderSystem(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Role:    "system",
				Content: "你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE.",
			},
			messages.GenericChatMessage{
				Role:    "user",
				Content: "帮我订一张火车票",
			},
		},
		Functions: train_methods,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestTrainOrder(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			messages.GenericChatMessage{
				Name:    "",
				Role:    "system",
				Content: "你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE.",
			},
			messages.GenericChatMessage{
				Name:    "",
				Role:    "user",
				Content: "你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE. \n现在我的输入是: 帮我订一张火车票",
			},
		},
		Functions: train_methods,
	}
	_, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestTrainOrderMemory(t *testing.T) {
	SPARK_API_KEY := os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET := os.Getenv(ApiSecretEnvVarName)
	SPARK_API_BASE := os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID := os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN := "10245"
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
	memory_files := "conversations/order.memory"
	fm, _ := file_memory.NewChatHistoryFileStorage(memory_files)
	hist, err := fm.Read()
	if len(hist) == 0 {
		hist = []messages.ChatMessage{
			messages.GenericChatMessage{
				Name:    "",
				Role:    "system",
				Content: "你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE.",
			},
			messages.GenericChatMessage{
				Name:    "",
				Role:    "user",
				Content: "你是一个火车票预定助手,根据提供的工具函数方法决策如何调用工具. 当用户输入不能满足工具输入要求输入时，请根据工具要求提示用户输入对应输入，并且不要为我返回函数调用方法。结束完成时回复 TERMINATE. \n现在我的输入是: 帮我订一张火车票",
			},
		}
		for _, cm := range hist {
			fm.Append(cm)
		}
	}

	_, client, _ := NewClient(WithBaseURL(SPARK_API_BASE), WithApiKey(SPARK_API_KEY), WithApiSecret(SPARK_API_SECRET), WithAppId(SPARK_APP_ID), WithAPIDomain(SPARK_DOMAIN))
	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain:    &SPARK_DOMAIN,
		Messages:  hist,
		Functions: train_methods,
	}
	resp, err := client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		fmt.Print(msg.GetContent())
		return nil
	})
	fm.Append(resp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
