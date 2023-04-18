/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  error.go  error.go 2022-11-30
 */

package constant

import (
	"github.com/lshaofan/go-framework/application/dto/response"
	"net/http"
)

var (
	ServerError = response.NewError(500, "服务器错误", nil, http.StatusInternalServerError)
	// UNAUTHORIZED 未登录
	UNAUTHORIZED = response.NewError(10000, "未登录", nil, http.StatusUnauthorized)
	// InvalidToken 非法Token
	InvalidToken = response.NewError(10001, "非法Token", nil, http.StatusUnauthorized)
	// TokenExpired Token过期
	TokenExpired = response.NewError(10002, "Token过期", nil, http.StatusUnauthorized)
	// UsernameOrPasswordError 用户名或密码错误
	UsernameOrPasswordError = response.NewError(10003, "用户名或密码错误", nil, http.StatusUnauthorized)
	// PlatformNotExist 平台不存在
	PlatformNotExist = response.NewError(10004, "平台不存在", nil, http.StatusPreconditionFailed)
	// PlatformIdCanNotEmpty 平台id不能为空
	PlatformIdCanNotEmpty = response.NewError(10005, "平台id不能为空", nil, http.StatusPreconditionFailed)
)
