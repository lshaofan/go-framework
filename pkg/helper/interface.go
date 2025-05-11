package helper

type IGlobal interface {
	// Init 初始化
	Init() error
}

// DefaultResultInterface 默认service返回数据结构
type DefaultResultInterface interface {
	IsError() bool
	GetError() *ErrorModel
	SetError(err error)
	GetData() any
	SetData(data any)
	SetResponse(data any, err error)
}
