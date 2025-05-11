package helper

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type ListRequest struct {
	Page     int    `form:"page" json:"page" query:"page" binding:"required"`
	PageSize int    `form:"page_size" json:"page_size" query:"page_size" binding:"required"`
	Order    string `form:"order" json:"order" query:"order" msg:"排序" `
	Field    string `form:"field" json:"field" query:"field" msg:"排序字段" `
}

type Validate struct {
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
}

var (
	validate     *Validate
	validateOnce sync.Once
)

func NewValidate() *Validate {
	validateOnce.Do(func() {
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

type Request struct {
	validateTags []string
}

// GetValidateErr 获取校验错误信息 传入错误对象和对象 对象tag为json form uri query header
func (r *Request) GetValidateErr(err error, obj interface{}) *ErrorModel {
	v := NewValidate()
	getObj := reflect.TypeOf(obj)
	var result *ErrorModel
	var errs validator.ValidationErrors
	// 判断err 是否是 validator.ValidationErrors 类型
	if !errors.As(err, &errs) {
		result = NewErrorModel(ERROR, err.Error(), nil, http.StatusPreconditionFailed)
	} else {

		for _, err := range errs {
			if f, exist := getObj.Elem().FieldByName(err.Field()); exist {
				var tag string
				for _, tagValueStr := range r.validateTags {
					tagStr, ok := f.Tag.Lookup(tagValueStr)
					if ok {
						tag = tagStr
						break
					}
				}
				result = NewErrorModel(
					ERROR,
					tag+strings.Replace(err.Translate(v.trans), err.Field(), "", -1),
					nil,
					http.StatusPreconditionFailed,
				)
			} else {
				result = NewErrorModel(
					ERROR,
					err.Translate(v.trans),
					nil,
					http.StatusPreconditionFailed,
				)
			}
			return result
		}

	}
	return result
}

func NewRequest() *Request {
	return &Request{
		validateTags: []string{"json", "form", "uri", "query", "header"},
	}
}
