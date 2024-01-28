package mq

import (
	"github.com/CloudStriver/cloudmind-mq/internal/svc"
	"github.com/CloudStriver/cloudmind-mq/internal/util/heap"
	"github.com/CloudStriver/cloudmind-mq/internal/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type CleanNotificationMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.CleanNotificationMessage
	Heap    *heap.PairHeap
}

func NewCleanNotificationMq(ctx context.Context, svcCtx *svc.ServiceContext) *CleanNotificationMq {
	CleanNotificationMq := &CleanNotificationMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		CleanNotificationMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	CleanNotificationMq.msgChan = make([]chan *message.CleanNotificationMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.CleanNotificationMessage, bufferCount)
		CleanNotificationMq.msgChan[i] = ch
		go CleanNotificationMq.consume(ch)
	}

	return CleanNotificationMq
}

func (l *CleanNotificationMq) Consume(_, value string) error {
	var msg []message.CleanNotificationMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CleanNotificationMq->Consume Unmarshal err : %v , val : %s", err, value)
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

func (l *CleanNotificationMq) consume(ch chan *message.CleanNotificationMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("CleanNotificationMq->consume err : %v", ok)
			return
		}
		if _, err := l.svcCtx.CloudMindSystemRPC.CleanNotification(l.ctx, &system.CleanNotificationReq{
			UserId: msg.UserId,
		}); err != nil {
			logx.Errorf("CleanNotificationMq->consume err : %v", err)
		} else {
			logx.Infof("CleanNotificationMq->consume message : %v", msg)
		}
	}
}
