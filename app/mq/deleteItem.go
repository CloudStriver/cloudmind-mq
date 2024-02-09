package mq

import (
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/heap"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type DeleteItemMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.DeleteItemMessage
	Heap    *heap.PairHeap
}

func NewDeleteItemMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteItemMq {
	DeleteItemMq := &DeleteItemMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		DeleteItemMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	DeleteItemMq.msgChan = make([]chan *message.DeleteItemMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.DeleteItemMessage, bufferCount)
		DeleteItemMq.msgChan[i] = ch
		go DeleteItemMq.consume(ch)
	}

	return DeleteItemMq
}

func (l *DeleteItemMq) Consume(_, value string) error {
	var msg []message.DeleteItemMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("DeleteItemMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	for _, d := range msg {
		num := l.Heap.Pop()
		l.msgChan[num.Second] <- &d
		num.First++
		l.Heap.Push(num)
	}
	return nil
}

func (l *DeleteItemMq) consume(ch chan *message.DeleteItemMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("DeleteItemMq->consume err : %v", ok)
			return
		}
		if _, err := l.svcCtx.CloudMindContentRPC.DeleteItem(l.ctx, &content.DeleteItemReq{
			ItemId: msg.ItemId,
		}); err != nil {
			logx.Errorf("DeleteItemMq->consume err : %v", err)
		} else {
			logx.Infof("DeleteItemMq->consume message : %v", msg)
		}
	}
}
