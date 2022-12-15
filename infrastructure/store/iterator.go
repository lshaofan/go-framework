/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  iterator.go  iterator.go 2022-11-30
 */

package store

// Iterator 迭代器
type Iterator struct {
	Index int
	Items []interface{}
}

func NewIterator(data []interface{}) *Iterator {
	return &Iterator{
		Items: data,
	}
}

// HasNext 判断是否有下一个
func (i *Iterator) HasNext() bool {
	if i.Items == nil && len(i.Items) == 0 {
		return false
	}
	return i.Index < len(i.Items)
}

// Next 获取下一个
func (i *Iterator) Next() interface{} {
	if i.HasNext() {
		i.Index++
		return i.Items[i.Index-1]
	}
	return nil
}
