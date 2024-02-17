package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFeedBackMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFeedBackMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFeedBackMq {
	return &CreateFeedBackMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFeedBackMq) Consume(_, value string) error {
	var msg *message.CreateFeedBackMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateFeedBackMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindContentRPC.CreateFeedBack(l.ctx, &content.CreateFeedBackReq{
		FeedbackType: msg.FeedbackType,
		UserId:       msg.UserId,
		ItemId:       msg.ItemId,
	}); err != nil {
		logx.Errorf("CreateFeedBackMq->consume err : %v", err)
		return err
	} else {
		logx.Infof("CreateFeedBackMq->consume message : %v", msg)
	}
	return nil
}
