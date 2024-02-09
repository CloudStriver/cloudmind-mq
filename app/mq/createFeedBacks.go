package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/heap"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFeedBacksMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan []*message.CreateFeedBacksMessage
	Heap    *heap.PairHeap
}

func NewCreateFeedBacksMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFeedBacksMq {
	CreateFeedBacksMq := &CreateFeedBacksMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		CreateFeedBacksMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	CreateFeedBacksMq.msgChan = make([]chan []*message.CreateFeedBacksMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan []*message.CreateFeedBacksMessage, bufferCount)
		CreateFeedBacksMq.msgChan[i] = ch
		go CreateFeedBacksMq.consume(ch)
	}

	return CreateFeedBacksMq
}

func (l *CreateFeedBacksMq) Consume(_, value string) error {
	var msg []*message.CreateFeedBacksMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateFeedBacksMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	num := l.Heap.Pop()
	l.msgChan[num.Second] <- msg
	num.First += len(msg)
	l.Heap.Push(num)
	return nil
}

func (l *CreateFeedBacksMq) consume(ch chan []*message.CreateFeedBacksMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("CreateFeedBacksMq->consume err : %v", ok)
			return
		}

		FeedBacks := lo.Map[*message.CreateFeedBacksMessage, *content.FeedBack](msg, func(m *message.CreateFeedBacksMessage, _ int) *content.FeedBack {
			return m.FeedBack
		})
		if _, err := l.svcCtx.CloudMindContentRPC.CreateFeedBacks(l.ctx, &content.CreateFeedBacksReq{
			FeedBacks: FeedBacks,
		}); err != nil {
			logx.Errorf("CreateFeedBacksMq->consume err : %v", err)
		} else {
			logx.Infof("CreateFeedBacksMq->consume message : %v", msg)
		}
	}
}
