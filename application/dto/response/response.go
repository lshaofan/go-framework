/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  response.go  response.go 2022-11-30
 */

package response

import (
	"github.com/gin-gonic/gin"
	"github.com/lshaofan/go-framework/application/dto/constant"
	"net/http"
)

// ErrorModel 错误模型
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

// Response 返回数据
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

// CreateOk 创建成功
func CreateOk(c *gin.Context) {
	Result(SUCCESS, nil, constant.CreateSuccess, http.StatusOK, c)
}

// CreateOkWithData 创建成功带返回数据
func CreateOkWithData(data, c *gin.Context) {
	Result(SUCCESS, data, constant.CreateSuccess, http.StatusOK, c)
}

// UpdateOk 更新成功
func UpdateOk(c *gin.Context) {
	Result(SUCCESS, nil, constant.UpdateSuccess, http.StatusOK, c)
}

// UpdateOkWithData 更新成功带返回数据
func UpdateOkWithData(data, c *gin.Context) {
	Result(SUCCESS, data, constant.UpdateSuccess, http.StatusOK, c)
}

// DeleteOk 删除成功
func DeleteOk(c *gin.Context) {
	Result(SUCCESS, nil, constant.DeleteSuccess, http.StatusOK, c)
}

// DeleteOkWithData 删除成功带返回数据
func DeleteOkWithData(data, c *gin.Context) {
	Result(SUCCESS, data, constant.DeleteSuccess, http.StatusOK, c)
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
