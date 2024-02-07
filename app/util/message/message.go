package message

import (
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
