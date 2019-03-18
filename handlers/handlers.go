package handlers

import (
	"github.com/anvie/port-scanner"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/state"
	. "github.com/prazd/nodes_mon_bot/state"
	"sync"
	"time"
)

// TODO: Balance check
//func NodesAlive(address string, port int, c chan bool) {
//	ps := portscanner.NewPortScanner(address, 2*time.Second, 5)
//	c <- ps.IsOpen(port)
//}

func Worker(wg *sync.WaitGroup, addr string, port int, r *NodesState) {
	defer wg.Done()
	ps := portscanner.NewPortScanner(addr, 2*time.Second, 5)
	isAlive := ps.IsOpen(port)
	r.Set(addr, isAlive)
}

func IsAlive(curr string, configData config.Config) string {

	nodesState := state.New()

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
	}

	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], port, nodesState)
	}
	wg.Wait()

	var message string

	for address, status := range nodesState.Result {
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
