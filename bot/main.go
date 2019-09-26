package main

import (
	"github.com/prazd/nodes_mon_bot/shared"
	"github.com/prazd/nodes_mon_bot/shared/keyboard"
	"log"
	"os"
	"time"

	"github.com/prazd/nodes_mon_bot/shared/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go shared.CheckStoppedList(b)

	b.Handle("/start", func(m *tb.Message) {
		err := shared.CheckUser(m.Sender.ID)
		if err != nil {
			log.Println(err)
			b.Send(m.Sender, "Problems...")
			return
		}
		b.Send(m.Sender, "Hi!I can help you with nodes monitoring!", &tb.SendOptions{ParseMode: "Markdown"},
			&tb.ReplyMarkup{ResizeReplyKeyboard: true, ReplyKeyboard: keyboard.MainMenu})
	})

	// Main handlers
	b.Handle(tb.OnText, func(m *tb.Message) {

		currenciesList := []string{"ETH", "ETC", "BCH", "BTC", "LTC"}

		var message string

		for _, currency := range currenciesList {
			if m.Text == currency {
				message, err = shared.GetMessageOfNodesState(strings.ToLower(currency))
				if err != nil {
					b.Send(m.Sender, "Error...")
					return
				} else {
					b.Send(m.Sender, message)
					return
				}
			}
		}

		message = "Sorry, this command doesn't exist"

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

	b.Start()
}
