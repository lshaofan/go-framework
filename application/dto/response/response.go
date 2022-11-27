package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorModel struct {
	Code       int         `json:"code" example:"-1"`
	Message    string      `json:"message" example:"操作失败"`
	Error      interface{} `json:"result"`
	HttpStatus int         `json:"httpStatus" swaggerignore:"true"`
}

const (
	ERROR   = -1
	SUCCESS = 0
)

type Response struct {
	Code    int         `json:"code" example:"0"`
	Result  interface{} `json:"result"`
	Message string      `json:"message" example:"操作成功"`
}

// PageList  分页数据
type PageList[T interface{}] struct {
	Total    int64 `json:"total" example:"100"`
	Data     []T   `json:"data" `
	Page     int   `json:"page" example:"1"`
	PageSize int   `json:"page_size" example:"10"`
}

// NewError 创建错误
func NewError(code int, message string, result interface{}, httpStatus int) *ErrorModel {
	return &ErrorModel{code, message, result, httpStatus}
}

func Result(code int, result interface{}, message string, httpStatus int, c *gin.Context) {
	// 开始时间
	c.JSON(httpStatus, Response{
		code,
		result,
		message,
	})
}

// Ok 操作成功
func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", http.StatusOK, c)
}

// OkWithMessage 带消息的操作成功
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, http.StatusOK, c)
}

// Success 成功返回数据
func Success(result interface{}, message string, c *gin.Context) {
	Result(SUCCESS, result, message, http.StatusOK, c)
}

// Fail 默认错误
func Fail(httpStatus int, c *gin.Context) {
	//Result(ERROR, map[string]interface{}{}, "操作失败", httpStatus, c)
	ThrowError(NewError(ERROR, "操作失败", nil, httpStatus), c)
}

// FailWithMessage 带消息的错误
func FailWithMessage(message string, httpStatus int, c *gin.Context) {
	//Result(ERROR, map[string]interface{}{}, message, httpStatus, c)
	ThrowError(NewError(ERROR, message, nil, httpStatus), c)
}

// ThrowError 抛出已知错误
func ThrowError(err *ErrorModel, c *gin.Context) {
	//Result(err.Code, err.Error, err.Message, err.HttpStatus, c)
	c.AbortWithStatusJSON(err.HttpStatus, Response{
		err.Code,
		err.Error,
		err.Message,
	})

}

// ParamError 参数校验错误
func ParamError(message string, c *gin.Context) {
	err := NewError(ERROR, message, nil, http.StatusPreconditionFailed)
	ThrowError(err, c)
}
