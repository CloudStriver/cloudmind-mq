package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	CreateNotificationsConf kq.KqConf
	ReadNotificationsConf   kq.KqConf
	CleanNotificationConf   kq.KqConf
}
