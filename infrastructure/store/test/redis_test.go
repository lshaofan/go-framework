/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  redis_test.go  redis_test.go 2022-11-30
 */

package test

import (
	"fmt"
	store2 "github.com/lshaofan/go-framework/infrastructure/store"
	"github.com/lshaofan/go-framework/utils"
	"testing"
	"time"
)

var conf store2.RedisConfig

const (
	// MiniProgramServiceListCacheKey 小程序首页服务列表缓存key
	MiniProgramServiceListCacheKey = "mini_program_service_list_platform_id:%d"
)

func init() {
	conf = store2.RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
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

// 测试设置hash
func TestHSet(t *testing.T) {
	ret := store2.NewOperation(conf).HSet("testHSet", "name13", "张w22三里斯").Unwrap()
	if ret.(int64) != 1 {
		t.Error("设置失败：", ret)
	}
	t.Log("设置结果：", ret)
}

// 测试HMSet
func TestHMSet(t *testing.T) {
	dat := make([]string, 0)
	dat = append(dat, "name1", "张三")
	dat = append(dat, "name2", "李四")
	data := map[string]interface{}{
		"name122": "张三1",
		"name222": dat,
	}
	stringString, err := utils.MapStringInterfaceToStringString(data)
	if err != nil {
		t.Error("设置失败：", err)
	}

	ret := store2.NewOperation(conf).HMSet("testHMSet", stringString).Unwrap()
	if !ret.(bool) {
		t.Error("设置失败：", ret)
	}
	t.Log("设置结果：", ret)
}

// 测试获取服务列表
func TestGetServiceList(t *testing.T) {

	ret := store2.NewOperation(conf).Get(fmt.Sprintf(MiniProgramServiceListCacheKey, 1))
	t.Log(ret)
}
