package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	CreateNotificationConf    kq.KqConf
	ReadNotificationsConf     kq.KqConf
	DeleteNotificationsConf   kq.KqConf
	CreateItemConf            kq.KqConf
	UpdateItemConf            kq.KqConf
	DeleteItemConf            kq.KqConf
	CreateFeedBackConf        kq.KqConf
	DeleteFileRelationConf    kq.KqConf
	DeleteCommentRelationConf kq.KqConf
	EnvHeader                 string
	CommentBatchSize          int64
}
