package mq

import (
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type DeleteNotificationsMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationsMq {
	return &DeleteNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotificationsMq) Consume(_, value string) error {
	var msg message.DeleteNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("DeleteNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}
	if _, err := l.svcCtx.CloudMindSystemRPC.DeleteNotifications(l.ctx, &system.DeleteNotificationsReq{
		UserId:          msg.UserId,
		NotificationIds: msg.NotificationIds,
		OnlyType:        msg.OnlyType,
	}); err != nil {
		logx.Errorf("DeleteNotificationsMq->consume err : %v", err)
	} else {
		logx.Infof("DeleteNotificationsMq->consume message : %v", msg)
	}
	return nil
}
