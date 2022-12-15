/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  logger.go  logger.go 2022-11-30
 */

package services

import (
	"github.com/lshaofan/go-framework/infrastructure/logger"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(args interface{}) *Logger {
	return &Logger{logger: logger.New(args)}
}

// AddErrorLog 添加错误日志
func (l *Logger) AddErrorLog(fields map[string]interface{}) {
	l.logger.WithFields(fields).Error()
}

// AddInfoLog 添加信息日志
func (l *Logger) AddInfoLog(fields map[string]interface{}) {
	l.logger.WithFields(fields).Info()
}
