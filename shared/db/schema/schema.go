package schema

type User struct {
	Telegram_id  int  `json:"telegram_id"`
	Subscription bool `json:"subscription"`
}

type NodeInfo struct {
	Addresses []string `json:"addresses"`
	Stopped   []string `json:"stopped"`
	Currency  string   `json:"currency"`
}

type NodesApi struct {
	Currency string `json:"currency"`
	Endpoint string `json:"endpoint"`
}
