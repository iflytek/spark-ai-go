package file_memory

import (
	"github.com/iflytek/spark-ai-go/sparkai/llms"
	"os"
)
import (
	"encoding/json"
)

type ChatHistoryFileStorage struct {
	file *os.File
}

func NewChatHistoryFileStorage(filename string) (*ChatHistoryFileStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ChatHistoryFileStorage{file: file}, nil
}

func (ls *ChatHistoryFileStorage) Append(log llms.ChatMessage) error {
	logData, err := json.Marshal(log)
	if err != nil {
		return err
	}
	logData = append(logData, '\n') // 每条日志以换行符结尾
	_, err = ls.file.Write(logData)
	if err != nil {
		return err
	}
	return nil
}

func (ls *ChatHistoryFileStorage) Read() ([]llms.ChatMessage, error) {
	var logs []llms.ChatMessage
	file, err := os.Open(ls.file.Name())
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	for {
		var log llms.GenericChatMessage
		if err := decoder.Decode(&log); err != nil {
			break // 读取完所有日志或者发生错误退出循环
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (ls *ChatHistoryFileStorage) Close() error {
	return ls.file.Close()
}
