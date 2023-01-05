package credential

import (
	"context"
	"github.com/lshaofan/go-framework/infrastructure/store"
	"testing"
)

var redisConfig *store.RedisConfig

func init() {
	redisConfig = &store.RedisConfig{
		Host:     "127.0.0.1",
		Port:     63791,
		Prefix:   "wechat:",
		Password: "",
		Database: 0,
	}

}

// 测试获取access_token
func TestDefaultAccessToken_GetAccessToken(t *testing.T) {
	c := store.NewOperation(redisConfig)
	i := NewDefaultAccessToken("wx1750b18bf6f25c0f", "fd3528dd79a88f421cb6ceb8056fe9fc", CacheKeyMiniProgramPrefix, c)
	ctx := context.Background()
	token, err := i.GetAccessToken(ctx)
	if err != nil {
		return
	}
	t.Log(token)
}

// 测试获取jsapi_ticket
func TestDefaultJsApiTicket_GetJsApiTicket(t *testing.T) {
	c := store.NewOperation(redisConfig)
	IAccess_token := NewDefaultAccessToken("wx1750b18bf6f25c0f", "fd3528dd79a88f421cb6ceb8056fe9fc", CacheKeyMiniProgramPrefix, c)
	ctx := context.Background()
	token, err := IAccess_token.GetAccessToken(ctx)
	t.Log("token:", token)
	if err != nil {
		return
	}
	IJsTicket := NewDefaultJsTicket("wx1750b18bf6f25c0f", CacheKeyMiniProgramPrefix, c)
	ticket, err := IJsTicket.GetTicket(token)
	if err != nil {
		return
	}
	t.Log("ticket:", ticket)
}
