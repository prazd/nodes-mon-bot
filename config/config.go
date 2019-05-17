package config

type NodeInfo struct {
	Port      int      `json:"port"`
	Addresses []string `json:"addresses"`
	Currency string `json:"currency"`
}
