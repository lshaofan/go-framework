package config

import "github.com/lshaofan/go-framework/infrastructure/store"

type Config struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	Cache     store.Operation
}

func NewConfig(appId, appSecret string, c store.RedisConfig) *Config {
	return &Config{
		AppId:     appId,
		AppSecret: appSecret,
		Cache:     *store.NewOperation(c),
	}
}
