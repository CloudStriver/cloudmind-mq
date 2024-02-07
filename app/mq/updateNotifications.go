package mq

import (
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/heap"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type UpdateNotificationsMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.UpdateNotificationsMessage
	Heap    *heap.PairHeap
}

func NewUpdateNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationsMq {
	UpdateNotificationsMq := &UpdateNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		UpdateNotificationsMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	UpdateNotificationsMq.msgChan = make([]chan *message.UpdateNotificationsMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.UpdateNotificationsMessage, bufferCount)
		UpdateNotificationsMq.msgChan[i] = ch
		go UpdateNotificationsMq.consume(ch)
	}

	return UpdateNotificationsMq
}

func (l *UpdateNotificationsMq) Consume(_, value string) error {
	var msg []message.UpdateNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("UpdateNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
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

func (l *UpdateNotificationsMq) consume(ch chan *message.UpdateNotificationsMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("UpdateNotificationsMq->consume err : %v", ok)
			return
		}
		if _, err := l.svcCtx.CloudMindSystemRPC.UpdateNotifications(l.ctx, &system.UpdateNotificationsReq{
			OnlyNotificationIds: msg.OnlyNotificationIds,
			OnlyUserId:          msg.OnlyUserId,
			OnlyType:            msg.OnlyType,
			OnlyIsRead:          msg.OnlyIsRead,
			IsRead:              msg.IsRead,
		}); err != nil {
			logx.Errorf("UpdateNotificationsMq->consume err : %v", err)
		} else {
			logx.Infof("UpdateNotificationsMq->consume message : %v", msg)
		}
	}
}
