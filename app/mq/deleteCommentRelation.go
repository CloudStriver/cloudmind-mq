package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type DeleteCommentRelationMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentRelationMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentRelationMq {
	return &DeleteCommentRelationMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentRelationMq) Consume(_, value string) error {
	var msg *message.DeleteCommentRelationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateNotificationMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	_, err := l.svcCtx.PlatformRPC.DeleteNode(l.ctx, &platform.DeleteNodeReq{
		NodeId:   msg.FromId,
		NodeType: msg.FromType,
	})
	if err != nil {
		logx.Errorf("DeleteCommentRelationMq->Consume DeleteNode err : %v , val : %s", err, value)
		return err
	}

	for {
		var res *platform.GetCommentListResp
		res, err = l.svcCtx.PlatformRPC.GetCommentList(l.ctx, &platform.GetCommentListReq{
			FilterOptions: &platform.CommentFilterOptions{
				OnlySubjectId: lo.ToPtr(msg.FromId),
			},
			Pagination: &basic.PaginationOptions{
				Limit: lo.ToPtr(l.svcCtx.Config.CommentBatchSize),
			},
		})
		if err != nil {
			logx.Errorf("DeleteCommentRelationMq->Consume GetCommentList err : %v , val : %s", err, value)
			return err
		}
		if len(res.Comments) == 0 {
			break
		}

		err = mr.Finish(func() error {
			ids := lo.Map(res.Comments, func(val *platform.CommentInfo, _ int) string {
				return val.Id
			})
			_, err1 := l.svcCtx.PlatformRPC.DeleteCommentByIds(l.ctx, &platform.DeleteCommentByIdsReq{Ids: ids})
			if err1 != nil {
				logx.Errorf("DeleteCommentRelationMq->Consume DeleteCommentByIds err : %v , val : %s", err1, ids)
				return err1
			}
			return nil
		}, func() error {
			for _, val := range res.Comments {
				_, err2 := l.svcCtx.PlatformRPC.DeleteNode(l.ctx, &platform.DeleteNodeReq{
					NodeId:   val.Id,
					NodeType: msg.FromType,
				})
				if err2 != nil {
					logx.Errorf("DeleteCommentRelationMq->Consume DeleteNode err : %v , val : %s", err2, val)
					return err2
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
