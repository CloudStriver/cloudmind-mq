package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileRelationMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileRelationMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileRelationMq {
	return &DeleteFileRelationMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileRelationMq) Consume(_, value string) error {
	var msg *message.DeleteFileRelationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateNotificationMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	for _, v := range msg.FromIds {
		_, err := l.svcCtx.RelationRPC.DeleteNode(l.ctx, &relation.DeleteNodeReq{
			NodeId:   v,
			NodeType: msg.FromType,
		})
		if err != nil {
			logx.Errorf("DeleteFileRelationMq->Consume DeleteNode err : %v , val : %s", err, v)
			return err
		}
	}

	return nil
}
