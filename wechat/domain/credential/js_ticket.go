package credential

import (
	"fmt"
	"github.com/lshaofan/go-framework/application/dto/response"
	"github.com/lshaofan/go-framework/infrastructure/store"
	"sync"
	"time"
)

// DefaultJsTicket 默认获取js ticket方法
type DefaultJsTicket struct {
	appID           string
	cache           store.Operation
	jsAPITicketLock *sync.Mutex

	Prefix string
}

// NewDefaultJsTicket new
func NewDefaultJsTicket(appID string, prefix string, cache *store.Operation) JsTicketHandle {
	return &DefaultJsTicket{
		appID:           appID,
		cache:           *cache,
		jsAPITicketLock: new(sync.Mutex),
		Prefix:          prefix,
	}
}

// ResTicket 请求jsapi_ticket返回结果
type ResTicket struct {
	response.ErrorModel
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

// GetTicket 获取jsapi_ticket
func (js *DefaultJsTicket) GetTicket(accessToken string) (ticketStr string, err error) {
	// 先从cache中取
	key := fmt.Sprintf("%s_jsapi_ticket_%s", js.Prefix, js.appID)
	if val := js.cache.Get(key); val.Err == nil {
		ticketStr = val.StringResult
		return
	}

	js.jsAPITicketLock.Lock()
	defer js.jsAPITicketLock.Unlock()

	// 双检，防止重复从微信服务器获取
	if val := js.cache.Get(key); val.Err == nil {
		ticketStr = val.StringResult
		return
	}

	var ticket ResTicket
	ticket, err = GetTicketFromServer(accessToken)
	if err != nil {
		return
	}
	expires := ticket.ExpiresIn - 1500
	ret := js.cache.Set(key, ticket.Ticket, store.WithExpire(time.Second*time.Duration(expires)))
	if ret.Result != store.SetSuccess {
		err = fmt.Errorf("获取access_token设置cache失败")
		return

	}
	ticketStr = ticket.Ticket
	return
}
