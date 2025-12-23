package helper

import (
	"context"
)

// BaseRequest 所有请求的基础结构
type BaseRequest struct {
	Ctx context.Context `json:"-" form:"-"`
}

// SetContext 设置上下文
func (r *BaseRequest) SetContext(ctx context.Context) {
	r.Ctx = ctx
}

// GetContext 获取上下文
func (r *BaseRequest) GetContext() context.Context {
	return r.Ctx
}

// IBaseRequest 接口定义
type IBaseRequest interface {
	SetContext(ctx context.Context)
	GetContext() context.Context
}
