package request

import (
	"mime/multipart"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/lshaofan/go-framework/application/dto/response"
	"github.com/lshaofan/go-framework/infrastructure/dao"
)

type GinRequest struct {
	c *gin.Context
}

func NewGinRequest(c *gin.Context) *GinRequest {
	return &GinRequest{c: c}
}

type Validate struct {
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
}

var (
	validate *Validate
	once     sync.Once
)

func NewValidate() *Validate {
	once.Do(func() {
		validate = &Validate{}
		//注册翻译器
		zh_ := zh.New()
		uni := ut.New(zh_, zh_)
		trans, _ := uni.GetTranslator("zh")
		//获取gin的校验器
		val := binding.Validator.Engine().(*validator.Validate)
		//注册翻译器
		_ = zh_translations.RegisterDefaultTranslations(val, trans)
		validate.validate = val
		validate.uni = uni
		validate.trans = trans
	})
	return validate
}

// TranslateToString    tag 为结构体的tag定义的验证字段来源  如：json form uri query  返回第一条验证错误信息的字符串
func (g *GinRequest) TranslateToString(tag string, err error, obj interface{}) string {
	v := NewValidate()
	getObj := reflect.TypeOf(obj)
	var result string
	// 判断err 是否是 validator.ValidationErrors 类型
	if _, ok := err.(validator.ValidationErrors); !ok {
		result = err.Error()
	} else {

		errors := err.(validator.ValidationErrors)
		for _, err := range errors {
			if f, exist := getObj.Elem().FieldByName(err.Field()); exist && f.Tag.Get(tag) != "" {
				result = f.Tag.Get(tag) + strings.Replace(err.Translate(v.trans), err.Field(), "", -1)
			} else {
				result = err.Translate(v.trans)
			}
			return result
		}
	}
	return result
}

// TranslateToMap Translate 翻译错误信息  tag 为结构体的tag定义的验证字段来源  如：json form uri query 返回所有验证错误信息的map
func (g *GinRequest) TranslateToMap(tag string, err error, obj interface{}) map[string][]string {
	v := NewValidate()
	getObj := reflect.TypeOf(obj)
	var result = make(map[string][]string)
	errors := err.(validator.ValidationErrors)

	for _, err := range errors {
		if f, exist := getObj.Elem().FieldByName(err.Field()); exist {

			str := f.Tag.Get(tag) + strings.Replace(err.Translate(v.trans), err.Field(), "", -1)
			result[err.Field()] = append(result[err.Field()], str)
		} else {
			result[err.Field()] = append(result[err.Field()], err.Translate(v.trans))
		}
	}
	return result
}

// Translate TranslateToMap Translate 翻译错误信息
func (g *GinRequest) Translate(err error) map[string][]string {
	v := NewValidate()
	var result = make(map[string][]string)
	errors := err.(validator.ValidationErrors)
	for _, err := range errors {
		result[err.Field()] = append(result[err.Field()], err.Translate(v.trans))
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
