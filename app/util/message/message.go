package message

import (
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
)

type CreateNotificationsMessage struct {
	TargetUserId    string
	SourceUserId    string
	SourceContentId string
	TargetType      int64
	Type            int64
	Text            string
}
type UpdateNotificationsMessage struct {
	UserId   string
	OnlyType *int64
}

type UpdateItemMessage struct {
	ItemId   string   `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	IsHidden *bool    `protobuf:"varint,2,opt,name=isHidden,proto3,oneof" json:"isHidden,omitempty"`
	Labels   []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	Comment  *string  `protobuf:"bytes,5,opt,name=comment,proto3,oneof" json:"comment,omitempty"`
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
