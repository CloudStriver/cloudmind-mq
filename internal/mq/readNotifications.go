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

type ReadNotificationsMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.ReadNotificationsMessage
	Heap    *heap.PairHeap
}

func NewReadNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *ReadNotificationsMq {
	ReadNotificationsMq := &ReadNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		ReadNotificationsMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	ReadNotificationsMq.msgChan = make([]chan *message.ReadNotificationsMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.ReadNotificationsMessage, bufferCount)
		ReadNotificationsMq.msgChan[i] = ch
		go ReadNotificationsMq.consume(ch)
	}

	return ReadNotificationsMq
}

func (l *ReadNotificationsMq) Consume(_, value string) error {
	var msg []message.ReadNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("ReadNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
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

func (l *ReadNotificationsMq) consume(ch chan *message.ReadNotificationsMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("ReadNotificationsMq->consume err : %v", ok)
			return
		}
		if _, err := l.svcCtx.CloudMindSystemRPC.ReadNotifications(l.ctx, &system.ReadNotificationsReq{
			FilterOptions: msg.NotificationFilterOptions,
		}); err != nil {
			logx.Errorf("ReadNotificationsMq->consume err : %v", err)
		} else {
			logx.Infof("ReadNotificationsMq->consume message : %v", msg)
		}
	}
}
