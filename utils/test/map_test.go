/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  map_test.go  map_test.go 2022-12-13
 */

package test

import "testing"

// 判断map是否有key
func MapHasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

func TestMapHasKey(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	if !MapHasKey(m, "key1") {
		t.Error("map has key1")
	}
	if MapHasKey(m, "key3") {
		t.Error("map has not key3")
	}
}
