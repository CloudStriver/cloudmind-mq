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

type CreateNotificationMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNotificationMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationMq {
	return &CreateNotificationMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotificationMq) Consume(_, value string) error {
	var msg *message.CreateNotificationMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateNotificationMq->Consume Unmarshal err : %v , val : %s", err, value)
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
		logx.Errorf("CreateNotificationMq->consume err : %v", err)
	}
	return nil
}
