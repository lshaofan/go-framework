/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  bind.go  bind.go 2022-11-30
 */

package request

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lshaofan/go-framework/application/dto/response"
	"mime/multipart"
	"reflect"
)

const (
	// FileNotExist 文件不存在
	FileNotExist = "file not exist"
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
		fmt.Println("bing err: ", err)
		response.ParamError(getValidMsg(err, obj), c)
		return
	}
	return nil
}

// Bing BingHeader header参数验证器
func Bing(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBind(obj)
	if err != nil {
		response.ParamError(getValidMsg(err, obj), c)
		return
	}
	return nil
}

// BingFile file参数验证器
func BingFile(c *gin.Context, name ...string) (file *multipart.FileHeader, err error) {
	if len(name) > 0 {
		file, err = c.FormFile(name[0])
		if err != nil {
			response.ParamError(FileNotExist, c)
			return
		}
		return
	} else {
		file, err = c.FormFile("file")
		if err != nil {
			response.ParamError(FileNotExist, c)
			return
		}
		return
	}
}
