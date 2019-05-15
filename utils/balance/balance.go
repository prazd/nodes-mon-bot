package balance

import (
	"github.com/onrik/ethrpc"
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

	balances["https://mainnet.infura.io"] = infuraResponse.String()

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

	// infura balance
	infuraInstance := ethrpc.New("https://etc-geth.0xinfra.com")

	infuraResponse, err := infuraInstance.EthGetBalance(address, "latest")
	if err != nil {
		return nil, err
	}

	balances["https://etc-geth.0xinfra.com"] = infuraResponse.String()

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

//func BchBalance(address string, endpoints []string) (Balances, error){
//	balances := new(Balances)
//
//	insightBalance := ""
//}
