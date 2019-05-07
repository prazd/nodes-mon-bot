package utils

import (
	"github.com/anvie/port-scanner"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/state"
	"sync"
	"time"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/prazd/nodes_mon_bot/subscribe"
	"reflect"
)


func Worker(wg *sync.WaitGroup, addr string, port int, r *state.NodesState) {
	defer wg.Done()
	ps := portscanner.NewPortScanner(addr, 1*time.Second, 1)
	isAlive := ps.IsOpen(port)
	r.Set(addr, isAlive)
}

func GetHostInfo(curr string, configData config.Config) ([]string, int) {

	var addresses []string
	var port int

	switch curr {
	case "eth":
		addresses = configData.EthNodes.Addresses
		port = configData.EthNodes.Port
	case "etc":
		addresses = configData.EtcNodes.Addresses
		port = configData.EtcNodes.Port
	case "btc":
		addresses = configData.BtcNodes.Addresses
		port = configData.BtcNodes.Port
	case "bch":
		addresses = configData.BchNodes.Addresses
		port = configData.BchNodes.Port
	case "ltc":
		addresses = configData.LtcNodes.Addresses
		port = configData.LtcNodes.Port
	case "xlm":
		addresses = configData.XlmNodes.Addresses
		port = configData.XlmNodes.Port
	}

	return addresses, port
}

func GetMessage(result map[string]bool) string {
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

func IsAlive(curr string, configData config.Config) string {

	nodesState := state.New()

	addresses, port := GetHostInfo(curr, configData)

	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], port, nodesState)
	}
	wg.Wait()

	message := GetMessage(nodesState.Result)

	return message
}


// Subscribe logic

func StartSubscribe(currency string, configData config.Config, bot *tb.Bot, m *tb.Message, Subscription *subscribe.Subscription){
	key := make(chan bool)

	_, ok := Subscription.Info[m.Sender.ID]
	if ok {
		isSubscribed := Subscription.Info[m.Sender.ID].Eth.IsSubscribed
		if isSubscribed{
			bot.Send(m.Sender, "Already subscribed!")
		}
	}

	Subscription.Set(m.Sender.ID, key, currency)

	nodesState := state.New()

	addresses, port := GetHostInfo(currency, configData)

	go func() {
		for {
			select {
			case <-key:
				return
			default:
				var wg sync.WaitGroup

				for i := 0; i < len(addresses); i++ {
					wg.Add(1)
					go Worker(&wg, addresses[i], port, nodesState)
				}
				wg.Wait()

				for _, alive := range nodesState.Result{
					if !alive {
						message := GetMessage(nodesState.Result)
						bot.Send(m.Sender,"Subscribe message: \n" + message)
					}
				}

				bot.Send(m.Sender, "All "+currency+" nodes checked")
				time.Sleep(time.Second * 2)
			}
		}
	}()

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
