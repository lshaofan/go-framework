package helper

import (
	"context"
)

// BaseRequest 所有请求的基础结构
type BaseRequest struct {
	Client ClientType      `json:"-" form:"-"`
	Ctx    context.Context `json:"-" form:"-"`
}

// SetClient 设置客户端类型
func (r *BaseRequest) SetClient(client ClientType) {
	r.Client = client
}

// GetClient 获取客户端类型
func (r *BaseRequest) GetClient() ClientType {
	return r.Client
}

// SetContext 设置上下文
func (r *BaseRequest) SetContext(ctx context.Context) {
	r.Ctx = ctx
}

// GetContext 获取上下文
func (r *BaseRequest) GetContext() context.Context {
	return r.Ctx
}

// 接口定义
type IBaseRequest interface {
	SetClient(client ClientType)
	GetClient() ClientType
	SetContext(ctx context.Context)
	GetContext() context.Context
}
