package credential

import "context"

type IAccessToken interface {
	// GetAccessToken 获取access_token
	GetAccessToken(ctx context.Context) (accessToken string, err error)
}
