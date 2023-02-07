/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  utils.go  utils.go 2022-11-30
 */

package request

import (
	"github.com/lshaofan/go-framework/infrastructure/dao"
)

// NewPageReq 初始化分页请求参数 默认第一页 每页10条
func NewPageReq() *dao.PageRequest {
	return &dao.PageRequest{
		Page:     1,
		PageSize: 10,
		Where:    make(map[string]interface{}),
		OrWhere:  make(map[string]interface{}),
	}
}
