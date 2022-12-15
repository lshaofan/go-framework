/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  structs.go  structs.go 2022-12-13
 */

package utils

import (
	"encoding/json"
	"github.com/fatih/structs"
)

// ToMapStringString ToMapToMapStringString 将结构体转换map[string]string
func ToMapStringString(v interface{}) (ret map[string]string, err error) {
	ret = make(map[string]string)
	rets := structs.Map(v)
	for k, v := range rets {
		// 判断是否是string类型
		if _, ok := v.(string); ok {
			ret[k] = v.(string)
		} else {
			v, err = json.Marshal(v)
			if err != nil {
				continue
			}
			v = string(v.([]byte))
			ret[k] = v.(string)
		}
	}
	return
}
