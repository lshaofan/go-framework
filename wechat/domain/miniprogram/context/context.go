package context

import (
	"github.com/lshaofan/go-framework/wechat/domain/credential"
	"github.com/lshaofan/go-framework/wechat/domain/miniprogram/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.IAccessToken
}
