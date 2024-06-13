package gormx

import (
	"fmt"

	"gorm.io/gorm/logger"

	"webapi/core/logx"
)

type Writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

// Printf 格式化打印日志
func (w *Writer) Printf(message string, data ...interface{}) {
	if logx.GetLevel() == logx.LogLevelDebug {
		logx.Debug(fmt.Sprintf(message, data...))
	}
}
