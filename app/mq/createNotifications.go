package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

const chanCount = 10
const bufferCount = 1024

type CreateNotificationsMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationsMq {
	return &CreateNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotificationsMq) Consume(_, value string) error {
	var msg *message.CreateNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindSystemRPC.CreateNotifications(l.ctx, &system.CreateNotificationsReq{
		TargetUserId:    msg.TargetUserId,
		SourceUserId:    msg.SourceUserId,
		SourceContentId: msg.SourceContentId,
		TargetType:      msg.TargetType,
		Type:            msg.Type,
		Text:            msg.Text,
	}); err != nil {
		logx.Errorf("CreateNotificationsMq->consume err : %v", err)
	}
	return nil
}
