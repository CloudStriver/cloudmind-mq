package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	CreateNotificationsConf kq.KqConf
	UpdateNotificationsConf kq.KqConf
	DeleteNotificationsConf kq.KqConf
	CreateItemsConf         kq.KqConf
	UpdateItemConf          kq.KqConf
	DeleteItemConf          kq.KqConf
	CreateFeedBacksConf     kq.KqConf
	EnvHeader               string
}
