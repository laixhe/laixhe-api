package core

import (
	"github.com/gofiber/contrib/fiberzap/v2"
)

// OrmWriter 日志写入器
type OrmWriter struct {
	logger *fiberzap.LoggerConfig
}

// NewOrmWriter 构造日志写入器
func NewOrmWriter(logger *fiberzap.LoggerConfig) *OrmWriter {
	return &OrmWriter{logger: logger}
}

// Printf 格式化打印日志
func (writer *OrmWriter) Printf(message string, data ...interface{}) {
	writer.logger.Infof(message, data...)
}
