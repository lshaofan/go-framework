/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  log.go  log.go 2022-11-30
 */

package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var (
	Log *logrus.Entry
)

func init() {
	lg := New(map[string]interface{}{"name": "system", "path": "info"})
	fields := make(map[string]interface{})
	fields["name"] = "system"
	fields["time"] = time.Now()
	Log = lg.WithFields(fields)
}

func New(args interface{}) *logrus.Logger {
	// 获取当前时间
	now := time.Now()
	// 获取当前年月日
	year, month, day := now.Date()
	nowTime := fmt.Sprintf("%d-%d-%d", year, month, day)
	// 取出args中的name，path, 如果没有则使用默认值
	argsName, ok := args.(map[string]interface{})["name"]
	if !ok {
		argsName = "app"
	}
	argsPath, ok := args.(map[string]interface{})["path"]
	if !ok {
		argsPath = "logs"
	}

	path := fmt.Sprintf("./runtime/log/%s/%s/%s.log", nowTime, argsPath, argsName)
	lg := logrus.New()
	lg.SetReportCaller(true)
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(24*time.Hour),  //最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
		rotatelogs.WithRotationCount(7),            //设置7份 大于7份 或到了清理时间 开始清理
		rotatelogs.WithRotationSize(100*1024*1024), //设置100MB大小,当大于这个容量时，创建新的日志文件

	)
	lg.SetOutput(io.MultiWriter(os.Stdout, writer))

	if gin.Mode() == gin.DebugMode {
		// 设置日志颜色
		lg.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		lg.SetLevel(logrus.DebugLevel)
	} else if gin.Mode() == gin.ReleaseMode {
		lg.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		lg.SetLevel(logrus.InfoLevel)
	}
	return lg
}
