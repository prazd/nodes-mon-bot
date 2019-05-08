package config

type NodeInfo struct {
	Port      int      `json:"port"`
	Addresses []string `json:"addresses"`
}

type Config struct {
	EthNodes NodeInfo `json:"ethNodes"`
	EtcNodes NodeInfo `json:"etcNodes"`
	BtcNodes NodeInfo `json:"btcNodes"`
	BchNodes NodeInfo `json:"bchNodes"`
	LtcNodes NodeInfo `json:"ltcNodes"`
	XlmNodes NodeInfo `json:"xlmNodes"`
}
