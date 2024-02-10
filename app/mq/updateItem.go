package mq

import (
	"fmt"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/heap"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type UpdateItemMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.UpdateItemMessage
	Heap    *heap.PairHeap
}

func NewUpdateItemMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemMq {
	UpdateItemMq := &UpdateItemMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		UpdateItemMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	UpdateItemMq.msgChan = make([]chan *message.UpdateItemMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.UpdateItemMessage, bufferCount)
		UpdateItemMq.msgChan[i] = ch
		go UpdateItemMq.consume(ch)
	}

	return UpdateItemMq
}

func (l *UpdateItemMq) Consume(_, value string) error {
	var msg []message.UpdateItemMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("UpdateItemMq->Consume Unmarshal err : %v , val : %s", err, value)
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

func (l *UpdateItemMq) consume(ch chan *message.UpdateItemMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("UpdateItemMq->consume err : %v", ok)
			return
		}
		fmt.Println(metadata.FromIncomingContext(l.ctx))
		if _, err := l.svcCtx.CloudMindContentRPC.UpdateItem(l.ctx, &content.UpdateItemReq{
			ItemId:   msg.ItemId,
			IsHidden: msg.IsHidden,
			Labels:   msg.Labels,
			Comment:  msg.Comment,
		}); err != nil {
			logx.Errorf("UpdateItemMq->consume err : %v", err)
		} else {
			logx.Infof("UpdateItemMq->consume message : %v", msg)
		}
	}
}
