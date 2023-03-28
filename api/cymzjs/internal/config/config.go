package config

import (
	"git.zc0901.com/go/god/api"
	"git.zc0901.com/go/god/lib/store/cache"
)

type Config struct {
	api.ServerConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	AppId     string // 微信小程序AppID
	AppSecret string // 微信小程序AppSecret

	Cache cache.ClusterConf

	Mysql struct {
		Cymzjs string
	}

	QiniuZhuKe struct {
		AK string
		SK string
	}
}
