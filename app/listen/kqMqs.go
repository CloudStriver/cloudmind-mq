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
		kq.MustNewQueue(c.CreateNotificationConf, mq.NewCreateNotificationMq(ctx, svcContext)),
		kq.MustNewQueue(c.CreateItemConf, mq.NewCreateItemMq(ctx, svcContext)),
		kq.MustNewQueue(c.UpdateItemConf, mq.NewUpdateItemMq(ctx, svcContext)),
		kq.MustNewQueue(c.CreateFeedBackConf, mq.NewCreateFeedBackMq(ctx, svcContext)),
		kq.MustNewQueue(c.DeleteItemConf, mq.NewDeleteItemMq(ctx, svcContext)),
		kq.MustNewQueue(c.DeleteNotificationsConf, mq.NewDeleteNotificationsMq(ctx, svcContext)),
		kq.MustNewQueue(c.DeleteFileRelationConf, mq.NewDeleteFileRelationMq(ctx, svcContext)),
		kq.MustNewQueue(c.DeleteCommentRelationConf, mq.NewDeleteCommentRelationMq(ctx, svcContext)),
	}
}
