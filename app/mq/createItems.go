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

type CreateItemsMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan []*message.CreateItemsMessage
	Heap    *heap.PairHeap
}

func NewCreateItemsMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateItemsMq {
	CreateItemsMq := &CreateItemsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		CreateItemsMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	CreateItemsMq.msgChan = make([]chan []*message.CreateItemsMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan []*message.CreateItemsMessage, bufferCount)
		CreateItemsMq.msgChan[i] = ch
		go CreateItemsMq.consume(ch)
	}

	return CreateItemsMq
}

func (l *CreateItemsMq) Consume(_, value string) error {
	var msg []*message.CreateItemsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateItemsMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	num := l.Heap.Pop()
	l.msgChan[num.Second] <- msg
	num.First += len(msg)
	l.Heap.Push(num)
	return nil
}

func (l *CreateItemsMq) consume(ch chan []*message.CreateItemsMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("CreateItemsMq->consume err : %v", ok)
			return
		}

		Items := lo.Map[*message.CreateItemsMessage, *content.Item](msg, func(m *message.CreateItemsMessage, _ int) *content.Item {
			return m.Item
		})
		if _, err := l.svcCtx.CloudMindContentRPC.CreateItems(l.ctx, &content.CreateItemsReq{
			Items: Items,
		}); err != nil {
			logx.Errorf("CreateItemsMq->consume err : %v", err)
		} else {
			logx.Infof("CreateItemsMq->consume message : %v", msg)
		}
	}
}
