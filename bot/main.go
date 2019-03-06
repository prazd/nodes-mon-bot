package main

import (
	"log"
	"os"
	"strconv"
	"time"

	. "github.com/prazd/nods_mon/handlers"
	. "github.com/prazd/nods_mon/keyboard"
	. "github.com/prazd/nods_mon/nodesAddresses"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Println(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {

		whitelist := [2]string{os.Getenv("DevOne"), os.Getenv("DevTwo")} //
		if strconv.Itoa(m.Sender.ID) == whitelist[0] || strconv.Itoa(m.Sender.ID) == whitelist[1] {
			b.Send(m.Sender, "Hi!I can help you with nodes monitoring!", &tb.ReplyMarkup{
				InlineKeyboard: MainMenu,
			})

			b.Handle(&AllButton, func(c *tb.Callback) {

				ethNodeAlive := make(chan bool)
				etcNodeAlive := make(chan bool)
				btcNodeAlive := make(chan bool)
				ltcNodeAlive := make(chan bool)
				bchNodeAlive := make(chan bool)

				var ethText, etcText, btcText, bchText, ltcText string

				go NodesAlive(EthNodeAddress, EthNodePort, ethNodeAlive)
				go NodesAlive(EtcNodeAddress, EtcNodePort, etcNodeAlive)
				go NodesAlive(BtcNodeAddress, BtcNodePort, btcNodeAlive)
				go NodesAlive(BchNodeAddress, BchNodePort, bchNodeAlive)
				go NodesAlive(LtcNodeAddress, LtcNodePort, ltcNodeAlive)

				for i := 0; i < 5; i++ {
					select {
					case eth := <-ethNodeAlive:
						ethText = strconv.FormatBool(eth)
					case etc := <-etcNodeAlive:
						etcText = strconv.FormatBool(etc)
					case btc := <-btcNodeAlive:
						btcText = strconv.FormatBool(btc)
					case bch := <-bchNodeAlive:
						bchText = strconv.FormatBool(bch)
					case ltc := <-ltcNodeAlive:
						ltcText = strconv.FormatBool(ltc)
					}
				}

				messageText :=
					"Nodes Alive:\n" +
						"ETH - " + ethText + "\n" +
						"ETC - " + etcText + "\n" +
						"BTC - " + btcText + "\n" +
						"BCH - " + bchText + "\n" +
						"LTC - " + ltcText

				b.Edit(c.Message, messageText, &tb.ReplyMarkup{
					InlineKeyboard: MainMenu,
				})
				b.Respond(c, &tb.CallbackResponse{})
			})
		}
	})
	b.Start()
}
