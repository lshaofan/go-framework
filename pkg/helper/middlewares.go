package helper

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}

// GinException 异常处理中间件
func GinException() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				PrintStack()
				logger := NewLogger(map[string]interface{}{"name": "exception", "path": "http"})
				logFields := make(map[string]interface{})
				// 记录日志
				logFields["ip"] = c.ClientIP()
				logFields["error"] = err
				logFields["请求地址"] = c.Request.URL
				logFields["method"] = c.Request.Method
				// 行号
				logFields["line"], _, _, _ = runtime.Caller(1)
				logger.AddErrorLog(logFields)
				var message string
				// 判断err类型
				switch expr := err.(type) {
				case string:
					message = expr
				case error:
					message = expr.Error()
				default:
					message = fmt.Sprintf("%v", expr)
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, NewResponse(
					ERROR,
					message,
					nil,
				))
				return
			}
		}()
		c.Next()
	}
}
