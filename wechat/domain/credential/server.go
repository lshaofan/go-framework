package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lshaofan/go-framework/utils"
)

const (
	// AccessTokenURL 获取access_token的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// AccessTokenURL 企业微信获取access_token的接口
	workAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// CacheKeyOfficialAccountPrefix 微信公众号cache key前缀
	CacheKeyOfficialAccountPrefix = "officialaccount_"
	// CacheKeyMiniProgramPrefix 小程序cache key前缀
	CacheKeyMiniProgramPrefix = "miniprogram_"
	// CacheKeyWorkPrefix 企业微信cache key前缀
	CacheKeyWorkPrefix = "work_"
	// getTicketURL 获取ticket的url
	getTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// GetAccessTokenFromServer GetAccessToken 从服务器获取accessToken
func GetAccessTokenFromServer(ctx context.Context, url string) (accessToken Result, err error) {
	var body []byte
	body, err = utils.HTTPGetContext(ctx, url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return
	}
	if accessToken.Code != 0 {
		err = fmt.Errorf("GetAccessTokenFromServer error : errcode=%v , errmsg=%v", accessToken.Code, accessToken.Message)
		return
	}
	return
}

// GetTicketFromServer 从服务器中获取ticket
func GetTicketFromServer(accessToken string) (ticket ResTicket, err error) {
	var response []byte
	url := fmt.Sprintf(getTicketURL, accessToken)
	response, err = utils.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &ticket)
	if err != nil {
		return
	}
	if ticket.Code != 0 {
		err = fmt.Errorf("getTicket Error : errcode=%d , errmsg=%s", ticket.Code, ticket.Message)
		return
	}
	return
}
