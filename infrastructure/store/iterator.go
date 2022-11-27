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
