package miniprogram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lshaofan/go-framework/utils"
	"github.com/lshaofan/go-framework/wechat/domain/credential"
	"github.com/lshaofan/go-framework/wechat/domain/miniprogram/auth"
	"github.com/lshaofan/go-framework/wechat/domain/miniprogram/config"
	wechatCtx "github.com/lshaofan/go-framework/wechat/domain/miniprogram/context"
)

// MiniProgram 微信小程序相关API
type MiniProgram struct {
	ctx  *wechatCtx.Context
	Auth *auth.Auth
}

func GetMiniProgram(cfg *config.Config) *MiniProgram {
	return NewMiniProgram(cfg)
}

func NewMiniProgram(cfg *config.Config) *MiniProgram {
	defaultAkHandle := credential.NewDefaultAccessToken(cfg.AppId, cfg.AppSecret, credential.CacheKeyMiniProgramPrefix, &cfg.Cache)
	ctx := &wechatCtx.Context{
		IAccessToken: defaultAkHandle,
		Config:       cfg,
	}
	return &MiniProgram{
		ctx:  ctx,
		Auth: auth.NewAuth(ctx),
	}
}

// HTTPGet HTTPGet请求
func (m *MiniProgram) HTTPGet(url string, params map[string]interface{}) (resp []byte, err error) {
	// 拼接url
	return
}

// HTTPPost HTTPPost请求
func (m *MiniProgram) HTTPPost(url string, params map[string]interface{}) (resp []byte, err error) {
	// 判断参数中是否有baseApi
	ak, err := m.ctx.GetAccessToken(context.Background())
	if err != nil {
		return
	}
	// 拼接url

	url = fmt.Sprintf("%s?access_token=%s", url, ak)

	// 将参数转换为
	jsonData, err := json.Marshal(params)
	if err != nil {
		return
	}
	// 发送请求
	resp, err = utils.HTTPPost(url, string(jsonData))
	if err != nil {
		return nil, err
	}
	return
}

// HTTPPostJSON HTTPPostJSON请求
func (m *MiniProgram) HTTPPostJSON(url string, params interface{}) (resp []byte, err error) {
	// 判断参数中是否有baseApi
	ak, err := m.ctx.GetAccessToken(context.Background())
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s?access_token=%s", url, ak)
	// 发送请求
	resp, err = utils.PostJSON(url, params)
	if err != nil {
		return nil, err
	}
	return
}
