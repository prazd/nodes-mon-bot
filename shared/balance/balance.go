package balance

import (
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"os"
)

type Balances map[string]string

func GetFormatMessage(balances Balances) string {
	var message string
	for endpoint, balance := range balances {
		message += endpoint + ": " + balance + "\n"
	}
	return message
}

func GetEthBalance(address string, endpoints []string) (Balances, error) {

	balances := make(Balances)

	// infura balance
	infuraInstance := ethrpc.New("https://mainnet.infura.io")

	infuraResponse, err := infuraInstance.EthGetBalance(address, "latest")
	if err != nil {
		return nil, err
	}

	balances["mainnet.infura.io"] = infuraResponse.String()

	var instance *ethrpc.EthRPC

	// check other nodes
	for _, ip := range endpoints {
		instance = ethrpc.New("http://" + ip + ":8545")

		response, err := instance.EthGetBalance(address, "latest")
		if err != nil {
			return nil, err
		}

		balances[ip] = response.String()
	}

	return balances, nil
}

func GetEtcBalance(address string, endpoints []string) (Balances, error) {
	balances := make(Balances)

	etcInstance := ethrpc.New("https://etc-geth.0xinfra.com")

	response, err := etcInstance.EthGetBalance(address, "latest")
	if err != nil {
		return nil, err
	}

	balances["etc-geth.0xinfra.com"] = response.String()

	var instance *ethrpc.EthRPC

	// check other nodes
	for _, ip := range endpoints {

		instance = ethrpc.New("http://" + ip + ":8545")

		response, err := instance.EthGetBalance(address, "latest")
		if err != nil {
			return nil, err
		}

		balances[ip] = response.String()
	}

	return balances, nil
}

func GetBtcBalance(address string, endpoints []string) (Balances, error) {

	balances := make(Balances)

	response, err := req.Get("https://insight.bitpay.com/api/addr/" + address + "/balance")
	if err != nil {
		return nil, err
	}

	balances["insight.bitpay.com"] = response.String()

	var balanceReq string

	// check other nodes
	for _, ip := range endpoints {

		balanceReq = "http://" + ip + "/insight-api/addr/" + address + "/balance"

		res, err := req.Get(balanceReq)
		if err != nil {
			return nil, err
		}

		balances[ip] = res.String()
	}

	return balances, nil
}

func GetLtcBalance(address string, endpoints []string) (Balances, error) {

	balances := make(Balances)

	response, err := req.Get("https://insight.litecore.io/api/addr/" + address + "/balance")
	if err != nil {
		return nil, err
	}

	balances["insight.litecore.io"] = response.String()

	var balanceReq string

	// check other nodes
	for _, ip := range endpoints {

		balanceReq = "http://" + ip + ":3001/api/addr/" + address + "/balance"

		res, err := req.Get(balanceReq)
		if err != nil {
			return nil, err
		}

		balances[ip] = res.String()
	}

	return balances, nil
}

func GetBchBalance(address string, endpoints []string) (Balances, error) {
	balances := make(Balances)

	response, err := req.Get("https://blockdozer.com/api/addr/" + address + "/balance")
	if err != nil {
		return nil, err
	}

	balances["blockdozer.com"] = response.String()

	var balanceReq string

	// check other nodes
	for _, ip := range endpoints {

		balanceReq = "http://" + ip + ":3000/api/addr/" + address + "/balance"

		res, err := req.Get(balanceReq)
		if err != nil {
			return nil, err
		}

		balances[ip] = res.String()
	}

	return balances, nil
}

func GetXlmBalance(address string, endpoints []string) (Balances, error) {

	balances := make(Balances)

	type StellarResponse struct {
		Balances []struct {
			Balance             string `json:"balance"`
			Buying_liabilities  string `json:"buying_liabilities"`
			Selling_liabilities string `json:"selling_liabilities"`
			Asset_type          string `json:"asset_type"`
		}
	}

	var response StellarResponse

	stellarResponse, err := req.Get(os.Getenv("xlm-api") + address)
	if err != nil {
		return nil, err
	}

	err = stellarResponse.ToJSON(&response)
	if err != nil {
		return nil, err
	}

	var stellarBalanceString string

	for _, j := range response.Balances {
		if j.Asset_type == "native" {
			stellarBalanceString = j.Balance
		}
	}

	if stellarBalanceString == "" {
		stellarBalanceString = "0"
	}

	balances[os.Getenv("xlm-api")] = stellarBalanceString

	return balances, nil
}
