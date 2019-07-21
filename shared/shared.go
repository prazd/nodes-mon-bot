package shared

import (
	"reflect"
	"sync"
	"time"

	"log"
	"strconv"

	"github.com/anvie/port-scanner"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"github.com/prazd/nodes_mon_bot/shared/balance"
	"github.com/prazd/nodes_mon_bot/shared/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
)

type NodesStatus struct {
	sync.Mutex
	Result map[string]bool
}

func New() *NodesStatus {
	return &NodesStatus{
		Result: make(map[string]bool),
	}
}

func (ds *NodesStatus) Set(key string, value bool) {
	ds.Lock()
	defer ds.Unlock()
	ds.Result[key] = value
}

type NodesInfo struct {
	Status    *NodesStatus
	Port      int
	Addresses []string
}

func Worker(wg *sync.WaitGroup, addr string, port int, r *NodesStatus) {
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
	if err != nil {
		return nil, 0, err
	}

	port, err := db.GetPort(currency)
	if err != nil {
		return nil, 0, err
	}

	return addresses, port, nil
}

func GetMessageWithResults(result map[string]bool) string {
	var message string

	if len(result) == 0 {
		message = "Running nodes count: 0"
		return message
	} else if len(result) > 10 {
		var runningNodesCount int
		var stoppedNodesInfo string
		for address, status := range result {
			switch status {
			case true:
				runningNodesCount++
			case false:
				stoppedNodesInfo += address + "\n"
			}
		}

		if len(stoppedNodesInfo) == 0 {
			message += "\nRunning nodes count: " + strconv.Itoa(runningNodesCount)
			return message

		} else {
			message += "\nStopped nodes:" + stoppedNodesInfo
			message += "\nRunning nodes count: " + strconv.Itoa(runningNodesCount)
			return message
		}

	}
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

	nodesState := New()
	if currency == "xlm" {
		message := "api: " + os.Getenv("xlm-api")
		return message, nil

	} else if currency == "bch" {
		message := "api: " + os.Getenv("bch-api")
		return message, nil
	}

	addresses, port, err := GetHostInfo(currency)
	if err != nil {
		return "", err
	}

	RunWorkers(addresses, port, nodesState)

	message := GetMessageWithResults(nodesState.Result)

	return message, nil
}

func RunWorkers(addresses []string, port int, state *NodesStatus) {
	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], port, state)
	}
	wg.Wait()
}

func CheckStoppedList(bot *tb.Bot) {

	stoppedNodesCount := map[string]int{
		"eth": 0,
		"etc": 0,
		"btc": 0,
		"ltc": 0,
	}
	var message string

	for currency, _ := range stoppedNodesCount {
		stoppedNodes, err := db.GetStoppedList(currency)
		if err != nil {
			log.Fatal(err)
		}
		stoppedNodesCount[currency] = len(stoppedNodes)
	}

	for {
		for currency, count := range stoppedNodesCount {
			stoppedNodes, err := db.GetStoppedList(currency)
			if err != nil {
				log.Fatal(err)
			}
			if len(stoppedNodes) > count {
				ids := db.GetAllSubscribers()
				if ids == nil {
					continue
				}

				difference := len(stoppedNodes) - count
				if difference > 1 {
					message += "Stopped nodes:\n"
					for i := 1; i <= difference; i++ {
						message += stoppedNodes[len(stoppedNodes)-i] + "\n"
					}
				} else {
					message += "Stopped node: " + stoppedNodes[len(stoppedNodes)-1]
				}

				for i := 0; i < len(ids); i++ {
					bot.Send(&tb.User{ID: ids[i]}, "Subscribe message:\nCurrency: "+currency+"\n"+message)
				}
				stoppedNodesCount[currency] = len(stoppedNodes)
				message = ""
			} else {
				stoppedNodesCount[currency] = len(stoppedNodes)
			}
		}
		time.Sleep(time.Second * 5)
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
		balances, err = balance.GetBchBalance(address, endpoints)
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
func GetApiBalance(currency, address string) (string, error) {

	type StellarBalance struct {
		Balances []struct {
			Balance             string `json:"balance"`
			Buying_liabilities  string `json:"buying_liabilities"`
			Selling_liabilities string `json:"selling_liabilities"`
			Asset_type          string `json:"asset_type"`
		}
	}

	endpoint, err := db.GetApiEndpoint(currency)
	if err != nil {
		return "", err
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

	if Contains(currency, btc) {

		type BTC struct {
			Balance string `json:"balance"`
		}

		var btc BTC

		respBtcBalance, err := req.Get(endpoint + address)
		if err != nil {
			return "", err
		}

		err = respBtcBalance.ToJSON(&btc)
		if err != nil {
			return "", err
		}
		return btc.Balance, nil

	} else if Contains(currency, eth) {
		var ethClient = ethrpc.New(endpoint)

		respEthBalance, err := ethClient.EthGetBalance(address, "latest")
		if err != nil {
			return "", err
		}
		return respEthBalance.String(), nil

	} else {

		var stellarBalance StellarBalance

		respXlmBalance, err := req.Get(endpoint + address)
		if err != nil {
			return "", nil
		}

		err = respXlmBalance.ToJSON(&stellarBalance)
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
