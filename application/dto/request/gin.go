package request

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/lshaofan/go-framework/application/dto/response"
	"github.com/lshaofan/go-framework/infrastructure/dao"
	"mime/multipart"
	"reflect"
	"strings"
)

type GinRequest struct {
	c        *gin.Context
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinRequest(c *gin.Context) *GinRequest {
	//注册翻译器
	zh_ := zh.New()
	uni := ut.New(zh_, zh_)
	trans, _ := uni.GetTranslator("zh")
	//获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	return &GinRequest{
		c:        c,
		uni:      uni,
		trans:    trans,
		validate: validate,
	}
}

// TranslateToString    tag 为结构体的tag定义的验证字段来源  如：json form uri query  返回第一条验证错误信息的字符串
func (g *GinRequest) TranslateToString(tag string, err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	var result string
	errors := err.(validator.ValidationErrors)
	for _, err := range errors {
		//TODO 如果翻译为空，就返回参数错误 ，按时不需要后续出现为空err情况再说
		//if err.Translate(g.trans) == "" {
		//	result = "参数错误"
		//	return result
		//}
		if f, exist := getObj.Elem().FieldByName(err.Field()); exist && f.Tag.Get(tag) != "" {
			result = f.Tag.Get(tag) + strings.Replace(err.Translate(g.trans), err.Field(), "", -1)
		} else {
			result = err.Translate(g.trans)
		}
		return result
	}

	return result
}

// TranslateToMap Translate 翻译错误信息  tag 为结构体的tag定义的验证字段来源  如：json form uri query 返回所有验证错误信息的map
func (g *GinRequest) TranslateToMap(tag string, err error, obj interface{}) map[string][]string {
	getObj := reflect.TypeOf(obj)
	var result = make(map[string][]string)
	errors := err.(validator.ValidationErrors)

	for _, err := range errors {
		if f, exist := getObj.Elem().FieldByName(err.Field()); exist {

			str := f.Tag.Get(tag) + strings.Replace(err.Translate(g.trans), err.Field(), "", -1)
			result[err.Field()] = append(result[err.Field()], str)
		} else {
			result[err.Field()] = append(result[err.Field()], err.Translate(g.trans))
		}
	}
	return result
}

// Translate TranslateToMap Translate 翻译错误信息
func (g *GinRequest) Translate(err error) map[string][]string {
	var result = make(map[string][]string)
	errors := err.(validator.ValidationErrors)
	for _, err := range errors {
		result[err.Field()] = append(result[err.Field()], err.Translate(g.trans))
	}
	return result
}

// BingFile file参数验证器
func (g *GinRequest) BingFile(c *gin.Context, name ...string) (file *multipart.FileHeader, err error) {
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

// NewPageOptions 获取分页options
func (g *GinRequest) NewPageOptions() *dao.PageRequest {
	return dao.NewPageReq()
}
