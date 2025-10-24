package helper

import (
	"errors"
	"fmt"
	"net/http"
)

// Response  返回数据用于api接口
type Response struct {
	Code    int    `json:"code" `
	Result  any    `json:"result"`
	Message string `json:"message" `
}

// NewResponse 创建返回数据
func NewResponse(code int, message string, result any) *Response {
	return &Response{Code: code, Result: result, Message: message}
}

// PageList  分页数据
type PageList[T interface{}] struct {
	Total    int64 `json:"total" `
	Data     []T   `json:"data" `
	Page     int   `json:"page" `
	PageSize int   `json:"page_size" `
}

func NewPageList[T interface{}]() *PageList[T] {
	return &PageList[T]{}
}

// DefaultResult 默认的返回数据结构,用于services处理完业务逻辑后返回给controller的数据结构
type DefaultResult struct {
	Err  *ErrorModel `json:"err"`
	Data any         `json:"data"`
}

// NewDefaultResult 创建默认的返回数据结构
func NewDefaultResult() *DefaultResult {
	return &DefaultResult{
		Err: nil,
	}
}

// IsError 判断时候有错误
func (r *DefaultResult) IsError() bool {
	return r.Err != nil
}

// GetError 获取错误信息
func (r *DefaultResult) GetError() *ErrorModel {
	return r.Err
}

// SetError 设置错误信息
func (r *DefaultResult) SetError(err error) {
	fmt.Println("SetError", err)
	if err == nil {
		return
	}
	fmt.Println("SetError2", err.Error())
	// 判断是否为ErrorModel
	var errModel *ErrorModel
	if errors.As(err, &errModel) {
		r.Err = errModel
		return
	}
	r.Err = NewErrorModel(
		-1,
		err.Error(),
		nil,
		http.StatusInternalServerError,
	)
}

// GetData 获取数据
func (r *DefaultResult) GetData() any {
	return r.Data
}

// SetData 设置数据
func (r *DefaultResult) SetData(data any) {
	r.Data = data
}

// SetResponse 设置返回数据
func (r *DefaultResult) SetResponse(data any, err error) {
	r.SetData(data)
	r.SetError(err)
}
