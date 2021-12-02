package notification_service_interface

import "github.com/mypurecloud/platform-client-sdk-go/v56/platformclientv2"

type NotificationService interface {
	Authenticate(config *platformclientv2.Configuration)
	CreateChannel(notificationsApi *platformclientv2.NotificationsApi) (string, string)
	SubscribeToTopic(notificationsApi *platformclientv2.NotificationsApi, channelId string)
	Listen(connectUri string)
}