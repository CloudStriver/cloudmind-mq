package mq

import (
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type UpdateItemMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateItemMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemMq {
	return &UpdateItemMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateItemMq) Consume(_, value string) error {
	var msg *message.UpdateItemMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("UpdateItemMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindContentRPC.UpdateItem(l.ctx, &content.UpdateItemReq{
		ItemId:   msg.ItemId,
		IsHidden: msg.IsHidden,
		Labels:   msg.Labels,
	}); err != nil {
		logx.Errorf("UpdateItemMq->consume err : %v", err)
		return err
	} else {
		logx.Infof("UpdateItemMq->consume message : %v", msg)
	}
	return nil
}
