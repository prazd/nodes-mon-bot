package shared

import (
	"reflect"
	"sync"
	"time"

	"log"
	"strconv"

	"github.com/anvie/port-scanner"
	"github.com/prazd/nodes_mon_bot/shared/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

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

func Worker(wg *sync.WaitGroup, addr string, r *NodesStatus) {
	defer wg.Done()

	u, _ := url.Parse(addr)

	var (
		port int
		host = u.Host
	)

	switch strings.Contains(host, ":8545") {
	case true:
		host = re.FindString(host)
		port = 8545
	default:
		port = 80
	}

	ps := portscanner.NewPortScanner(host, 3*time.Second, 1)
	isAlive := ps.IsOpen(port)
	if !isAlive {
		time.Sleep(time.Second * 5)
		secondCheck := ps.IsOpen(port)
		r.Set(addr, secondCheck)
		return
	}
	r.Set(addr, isAlive)
}

func RunWorkers(addresses []string, state *NodesStatus) {
	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], state)
	}
	wg.Wait()
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

	addresses, err := db.GetEndpointsByCurrency(currency)
	if err != nil {
		return "", err
	}

	RunWorkers(addresses, nodesState)

	return GetMessageWithResults(nodesState.Result), nil
}

func CheckStoppedList(bot *tb.Bot) {

	stoppedNodesCount := map[string]int{
		"eth": 0,
		"etc": 0,
		"btc": 0,
		"ltc": 0,
		"bch": 0,
	}
	var message string

	for currency := range stoppedNodesCount {
		stoppedNodes, err := db.GetStoppedList(currency)
		if err != nil {
			log.Fatal(err)
		}
		stoppedNodesCount[currency] = len(stoppedNodes)
	}

	log.Println("Start check!")

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
