package store

import (
	"context"
	"time"
)

// Operation redis 操作
type Operation struct {
	ctx context.Context
}

type RedisConfig struct {
	Host     string
	Port     int
	Prefix   string
	Password string
	Database int
}

var Conf *RedisConfig

func NewOperation(c *RedisConfig) *Operation {
	Conf = c
	return &Operation{
		ctx: context.Background(),
	}
}

// Get 获取单个值
func (o *Operation) Get(key string) *Result {
	return NewResult(Redis().Get(o.ctx, key).Result())
}

// MGet 获取多值
func (o *Operation) MGet(keys ...string) *Result {
	return NewResult(Redis().MGet(o.ctx, keys...).Result())
}

// Set 设置值
func (o *Operation) Set(key string, value interface{}, attrs ...*Attr) *InterfaceResult {
	key = Conf.Prefix + key
	exp := Attrs(attrs).FindAttr(AttrExpire)
	// setNx
	nx := Attrs(attrs).FindAttr(AttrNx).UnwrapWithDefault(nil)
	if nx != nil {
		return NewInterfaceResult(Redis().SetNX(o.ctx, key, value, exp.UnwrapWithDefault(time.Second*0).(time.Duration)).Result())
	}
	// setXx
	xx := Attrs(attrs).FindAttr(AttrXx).UnwrapWithDefault(nil)
	if xx != nil {
		return NewInterfaceResult(Redis().SetXX(o.ctx, key, value, exp.UnwrapWithDefault(time.Second*0).(time.Duration)).Result())
	}

	return NewInterfaceResult(Redis().Set(o.ctx, key, value, exp.UnwrapWithDefault(time.Second*0).(time.Duration)).Result())
}

// Del 删除值 返回删除的数量
func (o *Operation) Del(keys ...string) *InterfaceResult {
	for i, key := range keys {
		keys[i] = Conf.Prefix + key
	}
	return NewInterfaceResult(Redis().Del(o.ctx, keys...).Result())
}
