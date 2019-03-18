package config

type Config struct {
	EthNodes struct {
		Port      int      `json:"port"`
		Addresses []string `json:"addresses"`
	} `json:"ethNodes"`
	EtcNodes struct {
		Port      int      `json:"port"`
		Addresses []string `json:"addresses"`
	} `json:"etcNodes"`
	BtcNodes struct {
		Port      int      `json:"port"`
		Addresses []string `json:"addresses"`
	} `json:"btcNodes"`
	BchNodes struct {
		Port      int      `json:"port"`
		Addresses []string `json:"addresses"`
	} `json:"bchNodes"`
	LtcNodes struct {
		Port      int      `json:"port"`
		Addresses []string `json:"addresses"`
	} `json:"ltcNodes"`
}
