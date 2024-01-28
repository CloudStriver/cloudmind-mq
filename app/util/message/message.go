package message

import (
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
)

type CreateNotificationsMessage struct {
	Notification *system.Notification
}

type ReadNotificationsMessage struct {
	NotificationFilterOptions *system.NotificationFilterOptions
}

type CleanNotificationMessage struct {
	UserId string
}
