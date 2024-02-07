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

type DeleteNotificationsMq struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	msgChan []chan *message.DeleteNotificationsMessage
	Heap    *heap.PairHeap
}

func NewDeleteNotificationsMq(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationsMq {
	DeleteNotificationsMq := &DeleteNotificationsMq{
		ctx:    ctx,
		svcCtx: svcCtx,
		Heap:   &heap.PairHeap{},
	}
	for i := 0; i < 10; i++ {
		DeleteNotificationsMq.Heap.Push(heap.Pair{
			First:  0,
			Second: i,
		})
	}
	DeleteNotificationsMq.msgChan = make([]chan *message.DeleteNotificationsMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *message.DeleteNotificationsMessage, bufferCount)
		DeleteNotificationsMq.msgChan[i] = ch
		go DeleteNotificationsMq.consume(ch)
	}

	return DeleteNotificationsMq
}

func (l *DeleteNotificationsMq) Consume(_, value string) error {
	var msg []message.DeleteNotificationsMessage
	if err := sonic.Unmarshal(pconvertor.String2Bytes(value), &msg); err != nil {
		logx.Errorf("DeleteNotificationsMq->Consume Unmarshal err : %v , val : %s", err, value)
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

func (l *DeleteNotificationsMq) consume(ch chan *message.DeleteNotificationsMessage) {
	for {
		msg, ok := <-ch
		if !ok {
			logx.Errorf("DeleteNotificationsMq->consume err : %v", ok)
			return
		}
		if _, err := l.svcCtx.CloudMindSystemRPC.DeleteNotifications(l.ctx, &system.DeleteNotificationsReq{
			OnlyNotificationIds: msg.OnlyNotificationIds,
			OnlyUserId:          msg.OnlyUserId,
			OnlyType:            msg.OnlyType,
			OnlyIsRead:          msg.OnlyIsRead,
		}); err != nil {
			logx.Errorf("DeleteNotificationsMq->consume err : %v", err)
		} else {
			logx.Infof("DeleteNotificationsMq->consume message : %v", msg)
		}
	}
}
