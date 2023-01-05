package credential

import (
	"context"
	"fmt"
	"github.com/lshaofan/go-framework/application/dto/response"
	"github.com/lshaofan/go-framework/infrastructure/store"
	"sync"
	"time"
)

// DefaultAccessToken 默认AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cache           store.Operation
	accessTokenLock *sync.Mutex
	Prefix          string
}

// Result ResAccessToken struct
type Result struct {
	response.ErrorModel
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewDefaultAccessToken(appID string, appSecret string, Prefix string, cache *store.Operation) IAccessToken {
	return &DefaultAccessToken{
		appID:           appID,
		appSecret:       appSecret,
		cache:           *cache,
		accessTokenLock: new(sync.Mutex),
		Prefix:          Prefix,
	}
}

func (ak *DefaultAccessToken) GetAccessToken(ctx context.Context) (accessToken string, err error) {
	// 设置cache key
	key := fmt.Sprintf("%saccess_token_%s", ak.Prefix, ak.appID)
	// 从cache中获取
	if val := ak.cache.Get(key); val.Err == nil {
		accessToken = val.StringResult
		return accessToken, nil
	}
	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()
	// 从cache中获取
	if val := ak.cache.Get(key); val.Err == nil {
		accessToken = val.StringResult
		return accessToken, nil
	}
	// 请求微信服务器
	var result Result
	result, err = GetAccessTokenFromServer(ctx, fmt.Sprintf(accessTokenURL, ak.appID, ak.appSecret))
	if err != nil {
		return
	}
	// 设置时间-1500秒，是为了防止因为网络延迟等原因，导致token提前失效
	expires := result.ExpiresIn - 1500
	// 设置cache
	ret := ak.cache.Set(key, result.AccessToken, store.WithExpire(time.Second*time.Duration(expires)))
	if ret.Result != store.SetSuccess {
		err = fmt.Errorf("获取access_token设置cache失败")
		return
	}
	accessToken = result.AccessToken
	return
}
