package file_memory

import (
	"autospark-go/llms"
	"fmt"
	"testing"
)

func Test_Write(t *testing.T) {
	// 创建日志存储对象
	logStorage, err := NewChatHistoryFileStorage("logs.jsonl")
	if err != nil {
		fmt.Println("创建日志存储对象失败:", err)
		return
	}
	defer logStorage.Close()

	// 添加日志
	log1 := llms.GenericChatMessage{Role: "human", Content: "This is log 2"}
	err = logStorage.Append(log1)
	if err != nil {
		fmt.Println("添加日志失败:", err)
		return
	}

	log2 := llms.GenericChatMessage{Role: "ai", Content: "This is log 1"}
	err = logStorage.Append(log2)
	if err != nil {
		fmt.Println("添加日志失败:", err)
		return
	}

	// 读取日志
	logs, err := logStorage.Read()
	if err != nil {
		fmt.Println("读取日志失败:", err)
		return
	}
	fmt.Println("读取到的日志:")
	for _, log := range logs {
		fmt.Println(log.GetType(), log.GetContent())
	}
}
