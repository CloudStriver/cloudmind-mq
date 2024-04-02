package message

type CreateNotificationMessage struct {
	TargetUserId    string
	SourceUserId    string
	SourceContentId string
	TargetType      int64
	Type            int64
	Text            string
}

type UpdateItemMessage struct {
	ItemId   string   `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	IsHidden *bool    `protobuf:"varint,2,opt,name=isHidden,proto3,oneof" json:"isHidden,omitempty"`
	Labels   []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
}

type CreateFeedBackMessage struct {
	FeedbackType string `protobuf:"bytes,1,opt,name=feedbackType,proto3" json:"feedbackType,omitempty"`
	UserId       string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	ItemId       string `protobuf:"bytes,3,opt,name=itemId,proto3" json:"itemId,omitempty"`
}

type CreateItemMessage struct {
	ItemId   string   `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	IsHidden bool     `protobuf:"varint,2,opt,name=isHidden,proto3" json:"isHidden,omitempty"`
	Labels   []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	Category string   `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
}

type DeleteItemMessage struct {
	ItemId string
}

type DeleteNotificationsMessage struct {
	UserId          string   `json:"userId"`
	NotificationIds []string `json:"notificationIds"`
	OnlyType        *int64   `json:"onlyType"`
}

type DeleteFileRelationsMessage struct {
	FromType int64    `json:"fromType"`
	FromId   string   `json:"fromId"`
	ToType   int64    `json:"toType"`
	Files    []string `json:"files"`
}
