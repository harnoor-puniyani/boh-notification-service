package notifier

type NotificationChannel struct {
	Type    string `json:"type"`
	Contact string `json:"contact"`
}

type NotificationEvent struct {
	UserID              string                `json:"userId"`
	NotificationMessage string                `json:"notificationMessage"`
	Channels            []NotificationChannel `json:"channels"`
}
