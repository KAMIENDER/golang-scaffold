package logger

import (
	"context"

	"github.com/KAMIENDER/golang-scaffold/infra/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOG_PATH = "./log/"
)

type Logger struct {
	logger *zap.Logger
}

func getLogFileName() string {
	nowTime := common.GetNowTime()
	return common.GetFileDateFormat(nowTime) + ".txt"
}

func GetLogger(ctx context.Context) *Logger {
	// 创建一个zap配置对象
	config := zap.NewProductionConfig()
	// 设置日志级别为DEBUG
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	// 设置时间格式为ISO8601
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{LOG_PATH + getLogFileName()}

	// 创建一个zap.Logger对象
	logger := zap.Must(config.Build())
	return &Logger{
		logger: logger,
	}
}

func (g *Logger) Info(msg string, fields ...zapcore.Field) {
	defer g.logger.Sync()
	g.logger.Info(msg, fields...)
}

func (g *Logger) Error(msg string, fields ...zapcore.Field) {
	defer g.logger.Sync()
	g.logger.Error(msg, fields...)
}

func (g *Logger) Debug(msg string, fields ...zapcore.Field) {
	defer g.logger.Sync()
	g.logger.Debug(msg, fields...)
}
