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

type DeleteItemMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteItemMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteItemMq {
	return &DeleteItemMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteItemMq) Consume(_, value string) error {
	var msg message.DeleteItemMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("DeleteItemMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	if _, err := l.svcCtx.CloudMindContentRPC.DeleteItem(l.ctx, &content.DeleteItemReq{
		ItemId: msg.ItemId,
	}); err != nil {
		logx.Errorf("DeleteItemMq->consume err : %v", err)
	} else {
		logx.Infof("DeleteItemMq->consume message : %v", msg)
	}
	return nil
}
