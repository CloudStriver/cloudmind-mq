package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

	if err := mr.Finish(lo.Map(msg.Files, func(item *content.FileParameter, _ int) func() error {
		var i int64
		for i = 1; i <= l.svcCtx.Config.RelationLength; i++ {
			_, _ = l.svcCtx.RelationRPC.DeleteRelation(l.ctx, &relation.DeleteRelationReq{
				FromType:     msg.FromType,
				FromId:       msg.FromId,
				ToType:       msg.ToType,
				ToId:         item.FileId,
				RelationType: i,
			})
		}
		return nil
	})...); err != nil {
		return err
	}

	return nil
}
