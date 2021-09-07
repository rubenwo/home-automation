package notification

type Notification struct {
	Title string
	Body  string
}
type NotificationSubscriber struct {
	ClientPlatform         string `json:"client_platform"`
	FirebaseMessagingToken string `json:"firebase_messaging_token"`
}
