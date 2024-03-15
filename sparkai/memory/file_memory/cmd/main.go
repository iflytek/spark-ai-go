package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Log struct {
	Message string `json:"message"`
}

type LogStorage struct {
	file *os.File
}

func NewLogStorage(filename string) (*LogStorage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &LogStorage{file: file}, nil
}

func (ls *LogStorage) Append(log Log) error {
	_, err := ls.file.Seek(0, os.SEEK_END) // 重新定位文件指针到文件末尾
	if err != nil {
		return err
	}

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

func (ls *LogStorage) Read() ([]Log, error) {
	var logs []Log
	file, err := os.Open(ls.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		var log Log
		if err := decoder.Decode(&log); err != nil {
			if err.Error() == "EOF" {
				break // 读取完所有日志退出循环
			}
			return nil, err // 发生其他错误返回错误
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (ls *LogStorage) Close() error {
	return ls.file.Close()
}

func main() {
	// 创建日志存储对象
	logStorage, err := NewLogStorage("logs.jsonl")
	if err != nil {
		fmt.Println("创建日志存储对象失败:", err)
		return
	}
	defer logStorage.Close()

	// 添加日志
	log1 := Log{Message: "This is log 1"}
	err = logStorage.Append(log1)
	if err != nil {
		fmt.Println("添加日志失败:", err)
		return
	}

	log2 := Log{Message: "This is log 2"}
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
		fmt.Println(log.Message)
	}
}
