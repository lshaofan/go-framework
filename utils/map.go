/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  map.go  map.go 2022-12-13
 */

package utils

import (
	"encoding/json"
	"fmt"
)

const (
	// ErrMapKeyNotExist 不存在key
	ErrMapKeyNotExist = "%s：未配置"
)

// MapStringInterfaceToStringString map[string]interface{} 转换为 map[string]string
func MapStringInterfaceToStringString(m map[string]interface{}) (ret map[string]string, err error) {
	ret = make(map[string]string)
	for k, v := range m {
		// 转换成json字符串
		v, err = json.Marshal(v)
		if err != nil {
			continue
		}
		ret[k] = string(v.([]byte))
	}
	return
}

// GetMapValueByStringKey 将字符串转换为map，并根据指定的key获取值
func GetMapValueByStringKey(obj string, keys ...string) (ret map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(obj), &ret)
	if err != nil {
		return
	}
	for _, key := range keys {
		// 判断是否有key
		if !MapHasKey(ret, key) {
			err = fmt.Errorf(ErrMapKeyNotExist, key)
		}
	}
	return
}

// MapHasKey 判断map是否有key
func MapHasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}
