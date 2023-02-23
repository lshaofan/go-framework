/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  result.go  result.go 2022-11-30
 */

package store

import (
	"encoding/json"
	"github.com/lshaofan/go-framework/utils"
)

const (
	SetSuccess = "OK"
)

// Result redis get 结果
type Result struct {
	StringResult string
	Err          error
	SliceResult  []interface{}
	MapResult    map[string]interface{}
	MapStringStr map[string]string
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
	case map[string]interface{}:
		r.MapResult = expr
	case map[string]string:
		r.MapStringStr = expr
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

// UnwrapSlice 取结果切片
func (r *Result) UnwrapSlice() []interface{} {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.SliceResult
}

// UnwrapSliceWithDefault 取结果切片，如果有错误则返回默认值
func (r *Result) UnwrapSliceWithDefault(defaultValue []interface{}) []interface{} {
	if r.Err != nil {
		return defaultValue
	}
	return r.SliceResult
}

// UnwrapMap 取结果map
func (r *Result) UnwrapMap() map[string]interface{} {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.MapResult
}

// UnwrapMapWithDefault 取结果map，如果有错误则返回默认值
func (r *Result) UnwrapMapWithDefault(defaultValue map[string]interface{}) map[string]interface{} {
	if r.Err != nil {
		return defaultValue
	}
	return r.MapResult
}

// UnwrapMapStringStr 获取map string string
func (r *Result) UnwrapMapStringStr() map[string]string {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.MapStringStr
}

// UnwrapMapStringStrWithDefault 获取map string string 结果 如果有错误则返回默认值
func (r *Result) UnwrapMapStringStrWithDefault(defaultValue map[string]string) map[string]string {
	if r.Err != nil {
		return defaultValue
	}
	return r.MapStringStr
}

type dbGetter func() interface{}
type DBGetter struct {
	Getter dbGetter
	Err    error
	Attrs  []*Attr

	Operation *Operation
}

func NewDBGetter(getter dbGetter, operation *Operation, attrs ...*Attr) *DBGetter {
	return &DBGetter{Getter: getter, Attrs: attrs, Operation: operation}
}

// Get 取数据  取出数据后任何格式数据均需要json 反序列化
func (g *DBGetter) Get(key string) *Result {
	ret := g.Operation.Get(key).UnWarpWithDefault("")
	if ret == "" {
		data := g.Getter()
		// 判断类型是否是错误
		if err, ok := data.(error); ok {
			return NewResult("", err)
		}
		// json 序列化
		jsonData, err := json.Marshal(data)
		if err != nil {
			return NewResult("", err)
		}
		// 设置数据
		g.Operation.Set(key, string(jsonData))
		return NewResult(string(jsonData), nil)
	}
	return NewResult(ret, nil)
}

// GetHash 取数据Hash
func (g *DBGetter) GetHash(key string, field string) *Result {
	ret := g.Operation.HGet(key, field).UnWarpWithDefault("")
	if ret == "" {
		data := g.Getter()
		// 判断类型是否是错误
		if err, ok := data.(error); ok {
			return NewResult("", err)
		}
		// json 序列化
		jsonData, err := json.Marshal(data)
		if err != nil {
			return NewResult("", err)
		}
		// 设置数据
		g.Operation.HSet(key, field, string(jsonData))
		return NewResult(string(jsonData), nil)
	}
	return NewResult(ret, nil)
}

// GetHashAll 取数据Hash
func (g *DBGetter) GetHashAll(key string) *Result {
	ret := g.Operation.HGetAll(key).UnwrapMapWithDefault(map[string]interface{}{})
	if len(ret) == 0 {
		data := g.Getter()
		// 判断类型是否是错误
		if err, ok := data.(error); ok {
			return NewResult(nil, err)
		}
		// 判断类型是否是map
		if m, ok := data.(map[string]interface{}); ok {
			// 设置数据
			stringString, err := utils.MapStringInterfaceToStringString(m)
			if err != nil {
				return NewResult(nil, err)
			}
			g.Operation.HMSet(key, stringString)
			return NewResult(m, nil)
		}
		return NewResult(nil, nil)
	}
	return NewResult(ret, nil)
}

// GetMHash 取数据MHash
func (g *DBGetter) GetMHash(key string, fields ...string) *Result {
	ret := g.Operation.HMGet(key, fields...).UnwrapMapWithDefault(map[string]interface{}{})
	if len(ret) == 0 {
		data := g.Getter()
		// 判断类型是否是错误
		if err, ok := data.(error); ok {
			return NewResult(nil, err)
		}
		// 判断类型是否是map
		if m, ok := data.(map[string]interface{}); ok {
			// 设置数据
			stringString, err := utils.MapStringInterfaceToStringString(m)

			if err != nil {
				return NewResult(nil, err)
			}
			g.Operation.HMSet(key, stringString)
			return NewResult(m, nil)
		}
		return NewResult(nil, nil)
	}
	return NewResult(ret, nil)
}

// Iterator 迭代器
func (r *Result) Iterator() *Iterator {
	return NewIterator(r.SliceResult)
}

// InterfaceResult  set结果接口
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
