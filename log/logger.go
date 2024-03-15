package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	Logger = GetLogger()
}

func GetLogger() *zap.SugaredLogger {
	// 创建一个自定义的日志级别配置
	levelConfig := zap.NewAtomicLevelAt(zapcore.ErrorLevel) // 设置日志级别为WARN
	// 使用自定义的日志级别配置创建生产环境Logger
	productionConfig := zap.NewProductionConfig()
	productionConfig.Level = levelConfig // 设置Logger的日志级别
	// 构建Logger
	logger, err := productionConfig.Build()
	if err != nil {
		// 处理错误
		panic(err)
	}
	defer logger.Sync() // 确保在程序结束时日志被刷新
	return logger.Sugar()
	//defer agent.logger.Sync()
}
