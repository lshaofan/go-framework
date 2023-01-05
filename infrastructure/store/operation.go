/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  operation.go  operation.go 2022-11-30
 */

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
	key = Conf.Prefix + key
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

// HSet SetHash 设置hash值
func (o *Operation) HSet(key string, field string, value interface{}) *InterfaceResult {
	key = Conf.Prefix + key
	return NewInterfaceResult(Redis().HSet(o.ctx, key, field, value).Result())
}

// HMSet SetHashMulti 设置多个hash值
func (o *Operation) HMSet(key string, fields map[string]string) *InterfaceResult {
	key = Conf.Prefix + key
	return NewInterfaceResult(Redis().HMSet(o.ctx, key, fields).Result())
}

// HGet GetHash 获取hash值
func (o *Operation) HGet(key string, field string) *Result {
	key = Conf.Prefix + key
	return NewResult(Redis().HGet(o.ctx, key, field).Result())
}

// HGetAll GetHashAll 获取hash所有值
func (o *Operation) HGetAll(key string) *Result {
	key = Conf.Prefix + key
	return NewResult(Redis().HGetAll(o.ctx, key).Result())
}

// HMGet GetHashMulti 获取多个hash值
func (o *Operation) HMGet(key string, fields ...string) *Result {
	key = Conf.Prefix + key
	return NewResult(Redis().HMGet(o.ctx, key, fields...).Result())
}

// HDel DelHash 删除hash值
func (o *Operation) HDel(key string, fields ...string) *InterfaceResult {
	key = Conf.Prefix + key
	return NewInterfaceResult(Redis().HDel(o.ctx, key, fields...).Result())
}
