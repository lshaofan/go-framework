/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  middlewares.go  middlewares.go 2022-11-30
 */

package http

import (
	"github.com/gin-gonic/gin"
	res "github.com/lshaofan/go-framework/application/dto/response"
	repo "github.com/lshaofan/go-framework/domain/repository"
	appServices "github.com/lshaofan/go-framework/domain/services"
	"net/http"
)

// Exception 异常处理中间件
func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				logger := repo.ILogger(appServices.NewLogger(map[string]interface{}{"name": "exception", "path": "http"}))
				logFields := make(map[string]interface{})
				// 记录日志
				logFields["ip"] = c.ClientIP()
				logFields["error"] = err
				logFields["url"] = c.Request.URL
				logFields["method"] = c.Request.Method
				logger.AddErrorLog(logFields)
				c.AbortWithStatusJSON(http.StatusInternalServerError, res.Response{
					Code:    -1,
					Message: err.(error).Error(),
				})
				return
			}
		}()
		c.Next()
	}
}
