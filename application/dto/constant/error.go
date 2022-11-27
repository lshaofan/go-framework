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
)
