package test

import (
	store2 "github.com/lshaofan/go-framework/infrastructure/store"
	"testing"
	"time"
)

var conf *store2.RedisConfig

func init() {
	conf = &store2.RedisConfig{
		Host:     "127.0.0.1",
		Port:     63791,
		Prefix:   "gb_app:",
		Password: "",
		Database: 0,
	}
}
func TestRedisGet(t *testing.T) {
	ret := store2.NewOperation(conf).Get("name").UnWarp()
	t.Log(ret)

	ret = store2.NewOperation(conf).Get("no_name").UnWarpWithDefault("default")
	t.Log(ret)
}

func TestMGet(t *testing.T) {
	iter := store2.NewOperation(conf).MGet("name222", "name1", "name2").Iterator()

	for iter.HasNext() {
		t.Log(iter.Next())
		t.Log("index:", iter.Index)
	}
}

func TestSet(t *testing.T) {
	ret := store2.NewOperation(conf).Set("testSetNoWithExpire", "测试设置，不带过期时间")
	if ret.Result != store2.SetSuccess {
		t.Error("不设置过期时间设置失败：", ret.Err)
	}

	ret = store2.NewOperation(conf).Set("testSetWithExpire", "测试设置，带过期时间",
		store2.WithExpire(time.Second*10),
	)
	if ret.Result != store2.SetSuccess {
		t.Error("设置过期时间设置失败：", ret.Err)
	}
	//
	ret = store2.NewOperation(conf).Set("testSetWithExpire", "测试设置，带过期时间20秒",
		store2.WithExpire(time.Second*20),
		store2.WithNx(),
	)
	if !ret.Result.(bool) {
		t.Error("设置Nx过期时间设置失败：", ret.Err)
	}

	ret = store2.NewOperation(conf).Set("testSetXX", "测试设置，带过期时间20秒",
		store2.WithExpire(time.Second*20),
	)
	if ret.Result != store2.SetSuccess {
		t.Error("设置XX带过期时间设置失败：", ret.Err)
	}

	ret = store2.NewOperation(conf).Set("testSetXX", "测试设置，带过期时间20秒",
		store2.WithExpire(time.Second*20),
		store2.WithXx(),
	)
	if !ret.Result.(bool) {
		t.Error("设置SetXX过期设置失败：", ret.Err)
	}
}

func TestDel(t *testing.T) {
	ret := store2.NewOperation(conf).Del("testSetNoWithExpire", "testSetWithExpire", "testSetXX").UnwrapWithDefault(0)
	if ret.(int64) == 0 {
		t.Error("全部删除失败：", ret)
	}
	if ret.(int64) != 4 {
		t.Error("删除失败,正常删除条数：", ret)
	}

	t.Log("删除结果：", ret)
}
