package store

const (
	SetSuccess = "OK"
)

// Result redis get 结果
type Result struct {
	StringResult string
	Err          error
	SliceResult  []interface{}
}

func NewResult(result interface{}, err error) *Result {
	r := &Result{
		Err: err,
	}
	// 判断类型
	switch expr := result.(type) {
	case string:
		r.StringResult = expr
	case []interface{}:
		r.SliceResult = expr

	default:
		return r
	}
	return r
}

// UnWarp 获取结果
func (r *Result) UnWarp() string {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.StringResult

}

// UnWarpWithDefault 获取结果，如果有错误则返回默认值
func (r *Result) UnWarpWithDefault(defaultValue string) string {
	if r.Err != nil {
		return defaultValue
	}
	return r.StringResult
}

// Iterator 迭代器
func (r *Result) Iterator() *Iterator {
	return NewIterator(r.SliceResult)
}

// InterfaceResult  结果接口
type InterfaceResult struct {
	Result interface{}
	Err    error
}

func NewInterfaceResult(result interface{}, err error) *InterfaceResult {
	return &InterfaceResult{Result: result, Err: err}
}

func (i *InterfaceResult) Unwrap() interface{} {
	if i.Err != nil {
		panic(i.Err)
	}
	return i.Result
}

func (i *InterfaceResult) UnwrapWithDefault(defaultValue interface{}) interface{} {
	if i.Err != nil {
		return defaultValue
	}
	return i.Result

}
