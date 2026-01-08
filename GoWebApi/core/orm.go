package core

import (
	contribZap "github.com/gofiber/contrib/v3/zap"
)

// OrmWriter 日志写入器
type OrmWriter struct {
	logger *contribZap.LoggerConfig
}

// NewOrmWriter 构造日志写入器
func NewOrmWriter(logger *contribZap.LoggerConfig) *OrmWriter {
	return &OrmWriter{logger: logger}
}

// Printf 格式化打印日志
func (writer *OrmWriter) Printf(message string, data ...interface{}) {
	writer.logger.Infof(message, data...)
}
