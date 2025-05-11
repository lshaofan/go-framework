package helper

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// GetLoggerOutput 获取 SetOutput
func GetLoggerOutput(path, name string) (*rotatelogs.RotateLogs, error) {
	// 获取当前时间
	now := time.Now()
	// 获取当前年月日
	year, month, day := now.Date()
	nowTime := fmt.Sprintf("%d-%d-%d", year, month, day)
	filePath := fmt.Sprintf("./runtime/log/%s/%s/%s.log", nowTime, path, name)
	return rotatelogs.New(
		filePath+".%Y%m%d",
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithRotationTime(24*time.Hour),  //最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
		rotatelogs.WithRotationCount(7),            //设置7份 大于7份 或到了清理时间 开始清理
		rotatelogs.WithRotationSize(100*1024*1024), //设置100MB大小,当大于这个容量时，创建新的日志文件

	)
}
