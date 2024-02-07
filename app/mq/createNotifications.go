package mq

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/CloudStriver/cloudmind-mq/app/util/heap"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

const chanCount = 10
const bufferCount = 1024

type CreateNotificationsMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan []*message.CreateNotificationsMessage
	Heap    *heap.PairHeap
}

func NewCreateNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationsMq {
	CreateNotificationsMq := &CreateNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		CreateNotificationsMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	CreateNotificationsMq.msgChan = make([]chan []*message.CreateNotificationsMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan []*message.CreateNotificationsMessage, bufferCount)
		CreateNotificationsMq.msgChan[i] = ch
		go CreateNotificationsMq.consume(ch)
	}

	return CreateNotificationsMq
}

func (l *CreateNotificationsMq) Consume(_, value string) error {
	var msg []*message.CreateNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("CreateNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
		return err
	}

	num := l.Heap.Pop()
	l.msgChan[num.Second] <- msg
	num.First += len(msg)
	l.Heap.Push(num)
	return nil
}

func (l *CreateNotificationsMq) consume(ch chan []*message.CreateNotificationsMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("CreateNotificationsMq->consume err : %v", ok)
			return
		}

		notifications := lo.Map[*message.CreateNotificationsMessage, *system.NotificationInfo](msg, func(m *message.CreateNotificationsMessage, _ int) *system.NotificationInfo {
			return m.Notification
		})
		if _, err := l.svcCtx.CloudMindSystemRPC.CreateNotifications(l.ctx, &system.CreateNotificationsReq{
			Notifications: notifications,
		}); err != nil {
			logx.Errorf("CreateNotificationsMq->consume err : %v", err)
		} else {
			logx.Infof("CreateNotificationsMq->consume message : %v", msg)
		}
	}
}
