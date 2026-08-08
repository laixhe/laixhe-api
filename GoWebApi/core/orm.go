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
//
// GORM 的 SQL / 慢 SQL / 错误 SQL 均经此写入 zap debug 级日志;
// 是否逐条输出由 config.yaml 的 orm.log_level 决定 (缺省 4=Info 逐条, 设 3 只记录慢 SQL 与错误)。
func (writer *OrmWriter) Printf(message string, data ...interface{}) {
	writer.logger.Debugf(message, data...)
}
