package listen

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/config"
	"github.com/CloudStriver/cloudmind-mq/app/mq"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

// pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.CreateNotificationsConf, mq.NewCreateNotificationsMq(ctx, svcContext)),
		kq.MustNewQueue(c.DeleteNotificationsConf, mq.NewDeleteNotificationsMq(ctx, svcContext)),
		kq.MustNewQueue(c.UpdateNotificationsConf, mq.NewUpdateNotificationsMq(ctx, svcContext)),
	}
}
