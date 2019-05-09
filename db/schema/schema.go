package schema

type User struct {
	Telegram_id  int `json:"telegram_id"`
	Subscription bool `json:"subscription"`
}