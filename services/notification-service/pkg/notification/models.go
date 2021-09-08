package notification

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NotificationSubscriber struct {
	ClientPlatform         string `json:"client_platform"`
	FirebaseMessagingToken string `json:"firebase_messaging_token"`
}
