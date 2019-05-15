package main

import (
	"log"
	"os"
	"time"

	"encoding/json"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/keyboard"
	"github.com/prazd/nodes_mon_bot/utils"

	"github.com/prazd/nodes_mon_bot/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ReadConfig() (*config.Config, error) {

	defaultConfigPath, _ := filepath.Abs("../config/config.json")
	configFile, err := os.Open(defaultConfigPath)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var conf config.Config

	if err = json.Unmarshal(byteValue, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	configData, err := ReadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go utils.FullCheckOfNode(*configData, b)

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
		b.Send(m.Sender, utils.GetMessageOfNodesState("eth", *configData))
	})

	b.Handle(&keyboard.EtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.GetMessageOfNodesState("etc", *configData))
	})

	b.Handle(&keyboard.BtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.GetMessageOfNodesState("btc", *configData))
	})

	b.Handle(&keyboard.BchButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.GetMessageOfNodesState("bch", *configData))
	})

	b.Handle(&keyboard.LtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.GetMessageOfNodesState("ltc", *configData))
	})

	b.Handle(&keyboard.XlmButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.GetMessageOfNodesState("xlm", *configData))
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

		currency := params[1]
		address := params[2]

		message, err := utils.GetBalances(currency, address, *configData)
		if err != nil {
			b.Send(m.Sender, "Problems...")
			return
		}

		b.Send(m.Sender, message)
	})

	b.Start()
}
