package utils

import (
	"reflect"
	"sync"
	"time"

	"github.com/anvie/port-scanner"
	"github.com/prazd/nodes_mon_bot/db"
	"github.com/prazd/nodes_mon_bot/state"
	"github.com/prazd/nodes_mon_bot/utils/balance"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
)

type NodesInfo struct {
	State     *state.SingleState
	Port      int
	Addresses []string
}

func Worker(wg *sync.WaitGroup, addr string, port int, r *state.SingleState) {
	defer wg.Done()
	ps := portscanner.NewPortScanner(addr, 3*time.Second, 1)
	isAlive := ps.IsOpen(port)
	if !isAlive {
		time.Sleep(time.Second * 5)
		secondCheck := ps.IsOpen(port)
		r.Set(addr, secondCheck)
		return
	}
	r.Set(addr, isAlive)
}

func GetHostInfo(currency string) ([]string, int, error) {

	addresses, err := db.GetAddresses(currency)
	if err != nil{
		return nil, 0, err
	}

	port, err := db.GetPort(currency)
	if err != nil{
		return nil, 0, err
	}

	return addresses, port, nil
}

func GetMessageWithResults(result map[string]bool) string {
	var message string
	for address, status := range result {
		message += address
		switch status {
		case true:
			message += ": ✔"
		case false:
			message += ": ✖"
		}
		message += "\n"
	}

	return message
}

func GetMessageOfNodesState(currency string) (string, error) {

	nodesState := state.NewSingleState()

	addresses, port, err := GetHostInfo(currency)
	if err != nil{
		return "",err
	}

	RunWorkers(addresses, port, nodesState)

	message := GetMessageWithResults(nodesState.Result)

	return message, nil
}

func RunWorkers(addresses []string, port int, state *state.SingleState) {
	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], port, state)
	}
	wg.Wait()
}

func isAllNodesUp(addresses []string, port int, state *state.SingleState) bool {
	RunWorkers(addresses, port, state)
	for _, j := range state.Result {
		if j == false {
			return false
		}
	}
	return true
}

func GetAllNodesFromDB() (map[string]NodesInfo,error) {

	allEntrys,err := db.GetAllNodesEntrys()
	if err != nil{
		return nil, err
	}

	return map[string]NodesInfo{
		"ETH": NodesInfo{
			State:     state.NewSingleState(),
			Port:      allEntrys["eth"].Port,
			Addresses: allEntrys["eth"].Addresses,
		},
		"ETC": NodesInfo{
			State:     state.NewSingleState(),
			Port:      allEntrys["etc"].Port,
			Addresses: allEntrys["etc"].Addresses,
		},
		"BTC": NodesInfo{
			State:     state.NewSingleState(),
			Port:      allEntrys["btc"].Port,
			Addresses: allEntrys["btc"].Addresses,
		},
		"LTC": NodesInfo{
			State:      state.NewSingleState(),
			Port:       allEntrys["ltc"].Port,
			Addresses:  allEntrys["ltc"].Addresses,
		},
		"BCH": NodesInfo{
			State:     state.NewSingleState(),
			Port:       allEntrys["bch"].Port,
			Addresses:  allEntrys["bch"].Addresses,
		},
		"XLM": NodesInfo{
			State:     state.NewSingleState(),
			Port:       allEntrys["xlm"].Port,
			Addresses:  allEntrys["xlm"].Addresses,
		},
	}, nil
}

func FullCheckOfNode(bot *tb.Bot) {

	allNodes, err := GetAllNodesFromDB()
	if err != nil{
		log.Println(err)
	}

	for {
		for currency, nodesInfo := range allNodes {
			up := isAllNodesUp(nodesInfo.Addresses, nodesInfo.Port, nodesInfo.State)
			if !up {
				ids := db.GetAllSubscribers()
				if ids == nil {
					continue
				}
				message := GetMessageWithResults(nodesInfo.State.Result)
				for i := 0; i < len(ids); i++ {
					bot.Send(&tb.User{ID: ids[i]}, "Subscribe message:\nCurrency: "+currency+"\n"+message)
				}
			}
		}
		time.Sleep(time.Second * 60)
	}
}

func Contains(params ...interface{}) bool {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return true
		}
	}
	return false
}

func CheckUser(id int) error {
	inDb, err := db.IsInDb(id)
	if err != nil {
		return err
	}

	if !inDb {
		err = db.CreateUser(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBalances(currency string, address string) (string, error) {

	var balances balance.Balances

	switch currency {
		case "eth":
			endpoints, err := db.GetAddresses("eth")
			if err != nil {
				return "", err
			}

			balances, err = balance.GetEthBalance(address, endpoints)
			if err != nil {
				return "", err
			}

		case "etc":
			endpoints, err := db.GetAddresses("etc")
			if err != nil {
				return "", err
			}
			balances, err = balance.GetEtcBalance(address, endpoints)
			if err != nil {
				return "", err
			}

		case "btc":
			endpoints, err := db.GetAddresses("btc")
			if err != nil {
				return "", err
			}

			balances, err = balance.GetBtcBalance(address, endpoints)
			if err != nil {
				return "", err
			}

		case "ltc":
			endpoints, err := db.GetAddresses("ltc")
			if err != nil {
				return "", err
			}
			balances, err = balance.GetLtcBalance(address, endpoints)
			if err != nil {
				return "", err
			}

		case "bch":
			endpoints, err := db.GetAddresses("bch")
			if err != nil {
				return "", err
			}
			balances, err = balance.GetBchBalance(address,endpoints)
			if err != nil {
				return "", err
			}
		case "xlm":
			endpoints, err := db.GetAddresses("xlm")
			if err != nil {
				return "", err
			}
			balances, err = balance.GetXlmBalance(address, endpoints)
			if err != nil {
				return "", err
			}
		}

	result := balance.GetFormatMessage(balances)

	return result, nil
}

// for node api
// API balances
func GetApiBalance(currency, address string)(string , error){

	type StellarBalance struct {
		Balances []struct {
			Balance             string `json:"balance"`
			Buying_liabilities  string `json:"buying_liabilities"`
			Selling_liabilities string `json:"selling_liabilities"`
			Asset_type          string `json:"asset_type"`
		}
	}

	endpoint, err := db.GetApiEndpoint(currency)
	if err != nil{
		return "",err
	}

	btc := []string{
		"btc",
		"bch",
		"ltc",
	}
	eth := []string{
		"eth",
		"etc",
	}

	if Contains(currency, btc){
		balance, err := req.Get(endpoint + address)
		if err != nil{
			return "", err
		}
		return balance.String(), nil

	}else if Contains(currency, eth){
		var ethClient = ethrpc.New(endpoint)
		balance, err := ethClient.EthGetBalance(address,"latest")
		if err != nil{
			return "", err
		}
		return balance.String(), nil

	} else {

		var stellarBalance StellarBalance

		balance, err := req.Get(endpoint+address)
		if err != nil{
			return "", nil
		}

		err = balance.ToJSON(&stellarBalance)
		if err != nil {
			return "", err
		}

		var stellarBalanceString string

		for _, j := range stellarBalance.Balances {
			if j.Asset_type == "native" {
				stellarBalanceString = j.Balance
			}
		}

		if stellarBalanceString == "" {
			stellarBalanceString = "0"
		}

		return stellarBalanceString, nil

	}

}
