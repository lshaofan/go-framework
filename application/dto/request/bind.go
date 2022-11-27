package request

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lshaofan/go-framework/application/dto/response"
	"reflect"
)

// getValidMsg 参数验证器
func getValidMsg(err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				return f.Tag.Get("msg")
			}
		}
	}
	return err.Error()
}

// BingJson json参数验证器
func BingJson(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBindJSON(obj)
	if err != nil {
		response.ParamError(getValidMsg(err, obj), c)
		return
	}
	return nil
}

// BingQuery query参数验证器
func BingQuery(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBindQuery(obj)
	if err != nil {
		response.ParamError(getValidMsg(err, obj), c)
		return
	}
	return nil
}

// BingForm form参数验证器
func BingForm(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBind(obj)
	if err != nil {
		response.ParamError(getValidMsg(err, obj), c)
		c.Abort()
		return
	}
	return nil
}

// BingUri uri参数验证器
func BingUri(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBindUri(obj)
	if err != nil {
		response.ParamError(getValidMsg(err, obj), c)
		return
	}
	return nil
}
