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

type CreateItemMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateItemMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateItemMq {
	return &CreateItemMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateItemMq) Consume(_, value string) error {
	var msg *message.CreateItemMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateItemMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindContentRPC.CreateItem(l.ctx, &content.CreateItemReq{
		ItemId:   msg.ItemId,
		IsHidden: msg.IsHidden,
		Labels:   msg.Labels,
		Category: msg.Category,
	}); err != nil {
		logx.Errorf("CreateItemMq->consume err : %v", err)
		return err
	} else {
		logx.Infof("CreateItemMq->consume message : %v", msg)
	}
	return nil
}
