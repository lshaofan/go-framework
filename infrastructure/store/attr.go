/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  attr.go  attr.go 2022-11-30
 */

package store

import (
	"fmt"
	"time"
)

const (
	AttrExpire = "expire"
	AttrNx     = "nx"
	AttrXx     = "xx"
)

type empty struct{}
type Attr struct {
	Name  string
	Value interface{}
}

type Attrs []*Attr

// FindAttr 查找属性
func (a Attrs) FindAttr(name string) *InterfaceResult {
	for _, attr := range a {
		if attr.Name == name {
			return NewInterfaceResult(attr.Value, nil)
		}
	}
	return NewInterfaceResult(nil, fmt.Errorf("not found attr %s", name))
}

func WithExpire(t time.Duration) *Attr {
	return &Attr{
		Name:  AttrExpire,
		Value: t,
	}
}

func WithNx() *Attr {
	return &Attr{
		Name:  AttrNx,
		Value: empty{},
	}
}

func WithXx() *Attr {
	return &Attr{
		Name:  AttrXx,
		Value: empty{},
	}
}
