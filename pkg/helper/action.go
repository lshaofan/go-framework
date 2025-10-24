package helper

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BindOption func(obj any) error

type Action interface {
	// Success 成功并返回数据
	Success(data any)
	Error(err any)
	ThrowError(err *ErrorModel)
	ThrowValidateError(err error)
	Bind(param any, opts ...BindOption) error
	BindParam(param any) error
	BindUriParam(param any) error
	ShouldBindBodyWith(param any, bb binding.BindingBody) error
	ShouldBindWith(param any, bb binding.Binding) error
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
	_, err := a.PrepareRequest(req, ctxFunc, bindFuncs...)
	if err != nil {
		a.ThrowValidateError(err)
		return
	}

	// 调用服务
	result := serviceCall(req)
	a.HandleResult(result)
}

// ProcessCreate 处理创建请求
func (a *BaseAction) ProcessCreate(req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
}

// ProcessQuery 处理查询请求
func (a *BaseAction) ProcessQuery(req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
}

// ProcessUpdate 处理更新请求
func (a *BaseAction) ProcessUpdate(req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
}

// ProcessDelete 处理删除请求
func (a *BaseAction) ProcessDelete(req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
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

	// 设置客户端类型
	clientType := ClientType(a.Context.GetHeader(ClientHeaderKey))

	// 直接通过接口设置客户端类型
	req.SetClient(clientType)

	// 直接通过接口设置上下文
	req.SetContext(ctx)

	return ctx, nil
}

//==================================handler==================================

// HandleCreate 处理创建请求
func HandleCreate(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessCreate(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	})
}

// HandleList 处理列表请求
func HandleList(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessQuery(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	})
}

// HandleShow 处理详情请求
func HandleShow(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessQuery(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	}, func(i interface{}) error {
		return a.Action.BindUriParam(i)
	})
}

// HandleUpdate 处理更新请求
func HandleUpdate(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessUpdate(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	}, func(i interface{}) error {
		return a.Action.BindUriParam(i)
	})
}

// HandleEdit 处理编辑请求
func HandleEdit(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessUpdate(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	}, func(i interface{}) error {
		return a.Action.BindUriParam(i)
	})
}

// HandleDelete 处理删除请求
func HandleDelete(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc) {
	a := NewBaseAction(c)
	a.ProcessDelete(req, serviceCall, ctxFunc, func(i interface{}) error {
		return a.Action.BindParam(i)
	}, func(i interface{}) error {
		return a.Action.BindUriParam(i)
	})
}

// HandleCustom 处理自定义请求
func HandleCustom(c *gin.Context, req IBaseRequest, serviceCall func(IBaseRequest) *DefaultResult, ctxFunc BaseCtxFunc, bindFuncs ...func(interface{}) error) {
	a := NewBaseAction(c)
	a.Process(req, serviceCall, ctxFunc, bindFuncs...)
}
