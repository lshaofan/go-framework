package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinResponse struct {
	c   *gin.Context
	Msg struct {
		GetSuccess string
	}
}

func NewGinResponse(c *gin.Context) *GinResponse {
	g := &GinResponse{
		c: c,
	}
	g.Msg.GetSuccess = GetSuccess
	return g
}

// NewError 创建错误响应
func (g *GinResponse) NewError(code int, message string, result interface{}, httpStatus int) *ErrorModel {
	return &ErrorModel{code, message, result, httpStatus}
}

// ThrowError 抛出已知错误
func (g *GinResponse) ThrowError(err *ErrorModel, c *gin.Context) {
	c.AbortWithStatusJSON(err.HttpStatus, Response{
		err.Code,
		err.Result,
		err.Message,
	})

}

// FailWithMessage 带消息的错误
func (g *GinResponse) FailWithMessage(message string, httpStatus int, c *gin.Context) {
	ThrowError(NewError(ERROR, message, nil, httpStatus), c)
}

func (g *GinResponse) Result(code int, result interface{}, message string, httpStatus int, c *gin.Context) {
	// 开始时间
	c.JSON(httpStatus, Response{
		code,
		result,
		message,
	})
}

// Ok 操作成功
func (g *GinResponse) Ok(c *gin.Context) {
	g.Result(SUCCESS, map[string]interface{}{}, OkSuccess, http.StatusOK, c)
}

// OkWithMessage 带消息的操作成功
func (g *GinResponse) OkWithMessage(message string, c *gin.Context) {
	g.Result(SUCCESS, map[string]interface{}{}, message, http.StatusOK, c)
}

// Success 成功返回数据
func (g *GinResponse) Success(result interface{}, message string, c *gin.Context) {
	g.Result(SUCCESS, result, message, http.StatusOK, c)
}

// CreateOk 创建成功
func (g *GinResponse) CreateOk(c *gin.Context) {
	g.Result(SUCCESS, nil, CreateSuccess, http.StatusOK, c)
}

// CreateOkWithData 创建成功带返回数据
func (g *GinResponse) CreateOkWithData(data, c *gin.Context) {
	g.Result(SUCCESS, data, CreateSuccess, http.StatusOK, c)
}

// UpdateOk 更新成功
func (g *GinResponse) UpdateOk(c *gin.Context) {
	g.Result(SUCCESS, nil, UpdateSuccess, http.StatusOK, c)
}

// UpdateOkWithData 更新成功带返回数据
func (g *GinResponse) UpdateOkWithData(data, c *gin.Context) {
	g.Result(SUCCESS, data, UpdateSuccess, http.StatusOK, c)
}

// DeleteOk 删除成功
func (g *GinResponse) DeleteOk(c *gin.Context) {
	g.Result(SUCCESS, nil, DeleteSuccess, http.StatusOK, c)
}

// DeleteOkWithData 删除成功带返回数据
func (g *GinResponse) DeleteOkWithData(data, c *gin.Context) {
	g.Result(SUCCESS, data, DeleteSuccess, http.StatusOK, c)
}

// Fail 默认错误
func (g *GinResponse) Fail(httpStatus int, c *gin.Context) {
	//Result(ERROR, map[string]interface{}{}, "操作失败", httpStatus, c)
	g.ThrowError(g.NewError(ERROR, "操作失败", nil, httpStatus), c)
}

// ParamError 参数校验错误 传入错误信息
func (g *GinResponse) ParamError(message string, c *gin.Context) {
	err := g.NewError(ERROR, message, nil, http.StatusPreconditionFailed)
	g.ThrowError(err, c)
}

// ParamErrorWithData 参数校验错误 传入错误信息
func (g *GinResponse) ParamErrorWithData(message string, data interface{}, c *gin.Context) {
	err := g.NewError(ERROR, message, data, http.StatusPreconditionFailed)
	g.ThrowError(err, c)
}

// FailWithMessageByStatusBadRequest 400错误
func (g *GinResponse) FailWithMessageByStatusBadRequest(message string, c *gin.Context) {
	g.ThrowError(g.NewError(ERROR, message, nil, http.StatusBadRequest), c)
}

// SuccessWithData 成功返回数据
func (g *GinResponse) SuccessWithData(result interface{}, c *gin.Context) {
	g.Result(SUCCESS, result, GetSuccess, http.StatusOK, c)
}
