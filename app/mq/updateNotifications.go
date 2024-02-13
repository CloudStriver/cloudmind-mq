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

type UpdateNotificationsMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationsMq {
	return &UpdateNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotificationsMq) Consume(_, value string) error {
	var msg message.UpdateNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("UpdateNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindSystemRPC.UpdateNotifications(l.ctx, &system.UpdateNotificationsReq{
		UserId:   msg.UserId,
		OnlyType: msg.OnlyType,
	}); err != nil {
		logx.Errorf("UpdateNotificationsMq->consume err : %v", err)
		return err
	}
	return nil
}
