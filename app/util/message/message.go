package message

import (
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
)

type CreateNotificationsMessage struct {
	Notification *system.NotificationInfo
}
type UpdateNotificationsMessage struct {
	OnlyNotificationIds []string
	OnlyUserId          *string
	OnlyType            *int64
	OnlyIsRead          *bool
	IsRead              bool
}

type DeleteNotificationsMessage struct {
	OnlyNotificationIds []string
	OnlyUserId          *string
	OnlyType            *int64
	OnlyIsRead          *bool
}

type UpdateItemMessage struct {
	ItemId     string   `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	IsHidden   *bool    `protobuf:"varint,2,opt,name=isHidden,proto3,oneof" json:"isHidden,omitempty"`
	Labels     []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	Categories []string `protobuf:"bytes,4,rep,name=categories,proto3" json:"categories,omitempty"`
	Comment    *string  `protobuf:"bytes,5,opt,name=comment,proto3,oneof" json:"comment,omitempty"`
}

type CreateFeedBacksMessage struct {
	FeedBack *content.FeedBack
}

type CreateItemsMessage struct {
	Item *content.Item
}

type DeleteItemMessage struct {
	ItemId string
}
