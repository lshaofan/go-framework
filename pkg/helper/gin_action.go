package helper

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

// BindParam 绑定参数
func (g *GinActionImpl) BindParam(param interface{}) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	err := g.c.ShouldBind(param)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}
	return nil
}

// BindUriParam 绑定uri参数
func (g *GinActionImpl) BindUriParam(param interface{}) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	err := g.c.ShouldBindUri(param)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}

	return nil
}

// ShouldBindBodyWith 绑定body参数
func (g *GinActionImpl) ShouldBindBodyWith(param any, bb binding.BindingBody) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	err := g.c.ShouldBindBodyWith(param, bb)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}
	return nil
}

// ShouldBindWith 绑定参数
func (g *GinActionImpl) ShouldBindWith(param any, bb binding.Binding) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	err := g.c.ShouldBindWith(param, bb)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}
	return nil
}

// Bind 绑定参数 此方法需要传入多个gin的绑定方法
// 如果是绑定uri参数需要传入c.ShouldBindUri 需要将uri参数放在最后 否则会报错 结构体中的uri参数需要加上tag binding:"omitempty" 否则 c.ShouldBindUri 之前的绑定方法会报错
func (g *GinActionImpl) Bind(param any, opts ...BindOption) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	for _, opt := range opts {
		if err := opt(param); err != nil {
			return g.req.GetValidateErr(err, param)
		}
	}
	return nil
}

func NewGinActionImpl(c *gin.Context) *GinActionImpl {
	return &GinActionImpl{
		c:   c,
		req: NewRequest(),
	}
}
