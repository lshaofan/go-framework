package helper

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Action interface {
	// Success 成功并返回数据
	Success(data any)
	Error(err any)
	ThrowError(err *ErrorModel)
	ThrowValidateError(err error)
	BindParam(param any) error // 智能绑定（自动识别所有参数类型）
	CreateOK()
	UpdateOK()
	DeleteOK()
	SuccessWithMessage(message string, data any)
	CreateOkWithMessage(message string)
	UpdateOkWithMessage(message string)
	DeleteOkWithMessage(message string)
}

// BaseAction 提供所有 Action 的基础功能
type BaseAction struct {
	Action
	Context *gin.Context
}

type BaseCtxFunc func(ctx context.Context, c *gin.Context) context.Context

// NewBaseAction 创建基础 Action
func NewBaseAction(c *gin.Context) *BaseAction {
	return &BaseAction{
		Action:  NewGinActionImpl(c),
		Context: c,
	}
}

// BindParam 绑定请求参数（仅使用ShouldBind）
func (a *BaseAction) BindParam(req IBaseRequest, ctxFunc BaseCtxFunc) error {
	return a.Bind(req, ctxFunc, a.Context.ShouldBind)
}

// Bind 绑定请求参数
func (a *BaseAction) Bind(req IBaseRequest, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) error {
	_, err := a.PrepareRequest(req, ctxFunc, bindFuncs...)
	return err
}

// HandleResult 统一处理服务返回结果
func (a *BaseAction) HandleResult(result *DefaultResult) {
	if result.IsError() {
		errModel := result.GetError()
		a.ThrowError(errModel)
		return
	}
	a.Success(result.GetData())
}

// BindAndValidate 绑定并验证请求参数
func (a *BaseAction) BindAndValidate(req IBaseRequest, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) error {
	_, err := a.PrepareRequest(req, ctxFunc, bindFuncs...)
	return err
}

// Process 统一处理请求
func (a *BaseAction) Process(req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	ctx, err := a.PrepareRequest(req, ctxFunc, bindFuncs...)
	if err != nil {
		// ⚠️ 检查：如果响应已经被写入（ctxFunc 可能已返回响应），直接返回，不要再写入第二次
		if a.Context.Writer.Written() {
			return
		}

		a.ThrowValidateError(err)
		return
	}

	// 调用服务
	req.SetContext(ctx)
	result := serviceCall(req)

	a.HandleResult(result)
}

// PrepareRequest 统一处理请求参数绑定和上下文设置
func (a *BaseAction) PrepareRequest(req IBaseRequest, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) (context.Context, error) {
	// 绑定请求参数
	for _, bindFunc := range bindFuncs {
		if err := bindFunc(req); err != nil {
			return nil, err
		}
	}

	// 设置上下文
	ctx := ctxFunc(context.Background(), a.Context)

	// ⚠️ 检查：ctxFunc 中可能已经返回了响应（如权限检查失败）
	if a.Context.Writer.Written() {
		return nil, fmt.Errorf("ctxFunc已返回响应，状态码：%d", a.Context.Writer.Status())
	}

	// 直接通过接口设置上下文
	req.SetContext(ctx)

	return ctx, nil
}

//==================================handler==================================

// HandleRequest 统一处理所有类型的请求（智能绑定：自动识别 URI、Body、Query 等参数）
// 适用于所有 CRUD 操作：Create、List、Show、Update、Delete 等
// 这是推荐的标准处理方法，会自动：
//  1. 智能绑定参数（自动识别 URI、JSON、Query、Form 等）
//  2. 设置请求上下文
//  3. 调用服务层
//  4. 统一返回结果
func HandleRequest(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.Process(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i) // 智能绑定，自动识别所有参数类型
	})
}

// HandleCustom 处理自定义请求（支持自定义绑定函数）
func HandleCustom(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a := NewBaseAction(c)
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
}
