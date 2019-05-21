package main

import (
	"log"
	"os"
	"time"
	"github.com/prazd/nodes_mon_bot/keyboard"
	"github.com/prazd/nodes_mon_bot/utils"

	"github.com/prazd/nodes_mon_bot/db"
	tb "gopkg.in/tucnak/telebot.v2"

	"strings"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go utils.FullCheckOfNode(b)

	b.Handle("/start", func(m *tb.Message) {
		err := utils.CheckUser(m.Sender.ID)
		if err != nil {
			log.Println(err)
			b.Send(m.Sender, "Problems...")
			return
		}
		b.Send(m.Sender, "Hi!I can help you with nodes monitoring!", &tb.SendOptions{ParseMode: "Markdown"},
			&tb.ReplyMarkup{ResizeReplyKeyboard: true, ReplyKeyboard: keyboard.MainMenu})
	})

	// Main handlers
	b.Handle(&keyboard.EthButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("eth")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle(&keyboard.EtcButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("etc")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle(&keyboard.BtcButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("btc")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle(&keyboard.BchButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("bch")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle(&keyboard.LtcButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("ltc")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle(&keyboard.XlmButton, func(m *tb.Message) {
		message,err := utils.GetMessageOfNodesState("xlm")
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	// Subscribe handlers

	b.Handle(&keyboard.SubscriptionStatus, func(m *tb.Message) {
		message, err := db.GetSubStatus(m.Sender.ID)
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, message)
	})

	b.Handle("/sub", func(m *tb.Message) {
		err := db.SubscribeOrUnSubscribe(m.Sender.ID, true)
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, "Successfully **subscribed** on every currency!", &tb.SendOptions{ParseMode: "Markdown"})

	})

	b.Handle("/stop", func(m *tb.Message) {
		err := db.SubscribeOrUnSubscribe(m.Sender.ID, false)
		if err != nil {
			b.Send(m.Sender, "Please send /start firstly")
			return
		}
		b.Send(m.Sender, "Successfully **unsubscribed** on every currency!", &tb.SendOptions{ParseMode: "Markdown"})
	})

	// Balance handler

	b.Handle("/balance", func(m *tb.Message) {
		params := strings.Split(m.Text, " ")
		if len(params) < 3{
			b.Send(m.Sender, "Error!")
			return
		}

		if params[1] == "trust"{
			currency := params[2]
			address := params[3]
			message, err := utils.GetApiBalance(currency, address)
			if err != nil {
				b.Send(m.Sender, "Problems...")
				return
			}

			b.Send(m.Sender, message)
		}

		currency := params[1]
		address := params[2]

		message, err := utils.GetBalances(currency, address)
		if err != nil {
			b.Send(m.Sender, "Problems...")
			return
		}

		b.Send(m.Sender, message)
	})



	b.Start()
}
