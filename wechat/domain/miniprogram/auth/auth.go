package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lshaofan/go-framework/application/dto/response"
	"github.com/lshaofan/go-framework/utils"
	wechatCxt "github.com/lshaofan/go-framework/wechat/domain/miniprogram/context"
)

const (
	code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	checkEncryptedDataURL = "https://api.weixin.qq.com/wxa/business/checkencryptedmsg?access_token=%s"

	getPhoneNumber = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
)

// Auth 登录/用户信息
type Auth struct {
	*wechatCxt.Context
}

// NewAuth new auth
func NewAuth(ctx *wechatCxt.Context) *Auth {
	return &Auth{ctx}
}

// ResCode2Session 登录凭证校验的返回结果
type ResCode2Session struct {
	response.ErrorModel
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
}

// RspCheckEncryptedData .
type RspCheckEncryptedData struct {
	response.ErrorModel
	Vaild      bool `json:"vaild"`       // 是否是合法的数据
	CreateTime uint `json:"create_time"` // 加密数据生成的时间戳
}

// Code2Session 登录凭证校验。
func (auth *Auth) Code2Session(jsCode string) (result ResCode2Session, err error) {
	return auth.Code2SessionContext(context.Background(), jsCode)
}

// Code2SessionContext 登录凭证校验。
func (auth *Auth) Code2SessionContext(ctx context.Context, jsCode string) (result ResCode2Session, err error) {
	var res []byte
	if res, err = utils.HTTPGetContext(ctx, fmt.Sprintf(code2SessionURL, auth.AppId, auth.AppSecret, jsCode)); err != nil {
		return
	}
	if err = json.Unmarshal(res, &result); err != nil {
		return
	}
	if result.Code != 0 {
		err = fmt.Errorf("Code2Session error : errcode=%v , errmsg=%v", result.Code, result.Message)
		return
	}
	return
}
