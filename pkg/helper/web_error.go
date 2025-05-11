package helper

//  示例: ServerError = NewErrorModel(500, "服务器错误", nil, http.StatusInternalServerError)

// ErrorModel 错误模型
type ErrorModel struct {
	Code       int         `json:"code" `
	Message    string      `json:"message" `
	Result     interface{} `json:"result"`
	HttpStatus int         `json:"httpStatus" swaggerignore:"true"`
}

func NewErrorModel(code int, message string, result interface{}, httpStatus int) *ErrorModel {
	return &ErrorModel{Code: code, Message: message, Result: result, HttpStatus: httpStatus}
}

func (e *ErrorModel) Error() string {
	return e.Message
}
