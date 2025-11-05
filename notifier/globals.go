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

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AcsMessage struct {
	ChannelRegistrationId string   `json:"channelRegistrationId"`
	To                    []string `json:"to"`
	Kind                  string   `"json:"kind"`
	Content               string   `json:"content"`
}
