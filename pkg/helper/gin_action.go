package helper

import (
	"errors"
	"net/http"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 绑定策略类型
type bindStrategyType int

const (
	bindStrategyNormal  bindStrategyType = iota // 普通绑定（Body/Query/Form）
	bindStrategyUriOnly                         // 仅 URI 参数
	bindStrategyMixed                           // 混合参数（URI + Body/Query/Form）
)

// 绑定策略缓存（全局缓存，所有请求共享）
var (
	bindStrategyCache     = make(map[reflect.Type]bindStrategyType)
	bindStrategyCacheLock sync.RWMutex
)

type GinActionImpl struct {
	c   *gin.Context
	res *Response
	req *Request
}

/** =================================response================================= */

func (g *GinActionImpl) returnJsonWithStatusOK() {
	g.c.AbortWithStatusJSON(http.StatusOK, g.res)
}

func (g *GinActionImpl) returnJsonWithStatusBadRequest() {
	g.c.AbortWithStatusJSON(http.StatusBadRequest, g.res)
}

// ThrowError 抛出错误
func (g *GinActionImpl) ThrowError(err *ErrorModel) {
	g.c.AbortWithStatusJSON(err.HttpStatus, NewResponse(
		err.Code,
		err.Message,
		err.Result,
	))
}

// Error 失败
func (g *GinActionImpl) Error(err any) {
	g.res = NewResponse(ERROR, "", nil)
	// 判断err 类型
	switch err.(type) {
	case *ErrorModel:
		g.ThrowError(err.(*ErrorModel))
		return
	case string:
		g.res.Message = err.(string)
	case error:
		g.res.Message = err.(error).Error()

	default:
		g.res.Message = "未知错误"

	}
	g.returnJsonWithStatusBadRequest()
}

// ThrowValidateError 参数验证错误抛出异常
func (g *GinActionImpl) ThrowValidateError(err error) {
	//	判断是否为ErrorModel
	var errModel *ErrorModel
	if errors.As(err, &errModel) {
		g.ThrowError(errModel)
	} else {
		g.Error(err)
	}
}

// Success 成功
func (g *GinActionImpl) Success(data any) {
	g.res = NewResponse(SUCCESS, Succeed, data)
	g.returnJsonWithStatusOK()
}

// CreateOK 创建成功
func (g *GinActionImpl) CreateOK() {
	g.res = NewResponse(SUCCESS, CreateSuccess, nil)
	g.returnJsonWithStatusOK()
}

// UpdateOK 更新成功
func (g *GinActionImpl) UpdateOK() {
	g.res = NewResponse(SUCCESS, UpdateSuccess, nil)
	g.returnJsonWithStatusOK()
}

// DeleteOK 删除成功
func (g *GinActionImpl) DeleteOK() {
	g.res = NewResponse(SUCCESS, DeleteSuccess, nil)
	g.returnJsonWithStatusOK()
}

// SuccessWithMessage 成功并返回消息
func (g *GinActionImpl) SuccessWithMessage(message string, data interface{}) {
	g.res = NewResponse(SUCCESS, message, data)
	g.returnJsonWithStatusOK()
}

// CreateOkWithMessage 创建成功并返回消息
func (g *GinActionImpl) CreateOkWithMessage(message string) {
	g.res = NewResponse(SUCCESS, message, nil)
	g.returnJsonWithStatusOK()
}

// UpdateOkWithMessage 更新成功并返回消息
func (g *GinActionImpl) UpdateOkWithMessage(message string) {
	g.res = NewResponse(SUCCESS, message, nil)
	g.returnJsonWithStatusOK()
}

// DeleteOkWithMessage 删除成功并返回消息
func (g *GinActionImpl) DeleteOkWithMessage(message string) {
	g.res = NewResponse(SUCCESS, message, nil)
	g.returnJsonWithStatusOK()
}

/** =================================request================================= */

// BindParam 智能绑定参数（自动检测并处理 URI、Body、Query、Form 等参数）
// 会自动识别结构体中的标签类型，选择正确的绑定顺序，最后统一验证
func (g *GinActionImpl) BindParam(param interface{}) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}

	// 检测参数结构，决定绑定策略
	bindStrategy := g.detectBindStrategy(param)

	// 根据策略执行绑定
	var err error
	switch bindStrategy {
	case bindStrategyMixed:
		// 混合参数（URI + Body/Query/Form）
		err = g.bindMixedParams(param)
	case bindStrategyUriOnly:
		// 仅 URI 参数
		err = g.bindUriOnly(param)
	default:
		// 普通参数（Body/Query/Form）
		err = g.bindNormalParams(param)
	}

	return err
}

// detectBindStrategy 检测绑定策略（带缓存优化）
func (g *GinActionImpl) detectBindStrategy(param interface{}) bindStrategyType {
	paramType := reflect.TypeOf(param).Elem()

	// 1. 先尝试从缓存读取（快速路径，读锁）
	bindStrategyCacheLock.RLock()
	if strategy, exists := bindStrategyCache[paramType]; exists {
		bindStrategyCacheLock.RUnlock()
		return strategy
	}
	bindStrategyCacheLock.RUnlock()

	// 2. 缓存未命中，执行检测（慢速路径）
	strategy := g.detectBindStrategyInternal(paramType)

	// 3. 写入缓存（写锁）
	bindStrategyCacheLock.Lock()
	bindStrategyCache[paramType] = strategy
	bindStrategyCacheLock.Unlock()

	return strategy
}

// detectBindStrategyInternal 内部实现：实际检测逻辑
func (g *GinActionImpl) detectBindStrategyInternal(paramType reflect.Type) bindStrategyType {
	hasUri := false
	hasOther := false

	// 遍历字段检查标签
	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)

		// 检查是否有 uri 标签
		if _, ok := field.Tag.Lookup("uri"); ok {
			hasUri = true
		}

		// 检查是否有其他绑定标签（json、form、query、header 等）
		if _, ok := field.Tag.Lookup("json"); ok {
			hasOther = true
		}
		if _, ok := field.Tag.Lookup("form"); ok {
			hasOther = true
		}
		if _, ok := field.Tag.Lookup("query"); ok {
			hasOther = true
		}

		// 早期退出优化：如果已经确定是混合模式，无需继续检查
		if hasUri && hasOther {
			return bindStrategyMixed
		}
	}

	// 决定策略
	if hasUri && hasOther {
		return bindStrategyMixed // 混合参数
	} else if hasUri {
		return bindStrategyUriOnly // 仅 URI
	}
	return bindStrategyNormal // 普通参数
}

// bindMixedParams 绑定混合参数
func (g *GinActionImpl) bindMixedParams(param interface{}) error {
	// 1. 先绑定 URI 参数（不验证，不自动写入响应）
	// ⚠️ 不能使用 g.c.BindUri，因为它会自动写入 400 响应
	// 使用 mapUri 手动绑定 URI 参数，不触发验证和自动响应
	_ = g.mapUri(param)

	// 2. 再绑定 Body/Query/Form 参数（不验证）
	// 优先绑定 Body（如果存在）
	if g.c.Request.Body != nil && g.c.Request.ContentLength > 0 {
		// 使用 ShouldBindBodyWith 绑定 JSON Body（支持重复读取）
		// 忽略错误，因为可能没有 JSON 字段或 Body 不是 JSON 格式
		_ = g.c.ShouldBindBodyWith(param, binding.JSON)
	}

	// 绑定 Query 参数（忽略错误）
	_ = g.c.ShouldBindQuery(param)

	// 3. 最后统一验证
	if err := binding.Validator.ValidateStruct(param); err != nil {
		return g.req.GetValidateErr(err, param)
	}

	return nil
}

// mapUri 手动绑定 URI 参数（不验证，不自动写入响应）
// 这是一个安全的 URI 绑定方法，不会触发 Gin 的自动 400 响应
func (g *GinActionImpl) mapUri(param interface{}) error {
	// 将 gin.Context.Params 转换为 map[string][]string 格式
	uriParams := make(map[string][]string)
	for _, p := range g.c.Params {
		uriParams[p.Key] = []string{p.Value}
	}

	// 使用 binding.Uri 的 BindUri 方法绑定
	return binding.MapFormWithTag(param, uriParams, "uri")
}

// bindUriOnly 仅绑定 URI 参数
func (g *GinActionImpl) bindUriOnly(param interface{}) error {
	err := g.c.ShouldBindUri(param)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}
	return nil
}

// bindNormalParams 绑定普通参数（Body/Query/Form）
func (g *GinActionImpl) bindNormalParams(param interface{}) error {
	err := g.c.ShouldBind(param)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}
	return nil
}

func NewGinActionImpl(c *gin.Context) *GinActionImpl {
	return &GinActionImpl{
		c:   c,
		req: NewRequest(),
	}
}
