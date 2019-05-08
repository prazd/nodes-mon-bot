package main

import (
	"log"
	"os"
	"time"

	"encoding/json"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/keyboard"
	"github.com/prazd/nodes_mon_bot/subscription"
	"github.com/prazd/nodes_mon_bot/utils"
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

	// channels for subscribe
	Subscription := subscription.SubNew()

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!I can help you with nodes monitoring!", &tb.SendOptions{ParseMode: "Markdown"},
			&tb.ReplyMarkup{ResizeReplyKeyboard: true, ReplyKeyboard: keyboard.MainMenu})
	})

	// Main handlers
	b.Handle(&keyboard.EthButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("eth", *configData))
	})

	b.Handle(&keyboard.EtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("etc", *configData))
	})

	b.Handle(&keyboard.BtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("btc", *configData))
	})

	b.Handle(&keyboard.BchButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("bch", *configData))
	})

	b.Handle(&keyboard.LtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("ltc", *configData))
	})

	b.Handle(&keyboard.XlmButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("xlm", *configData))
	})

	// Subscribe handlers

	b.Handle(&keyboard.SubscriptionStatus, func(m *tb.Message) {
		b.Send(m.Sender, utils.SubStatus(&Subscription, m.Sender.ID))
	})

	b.Handle("/sub", func(m *tb.Message) {

		params := strings.Split(m.Text, " ")

		switch params[1] {
		case "eth":
			utils.StartSubscribe("eth", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Eth.IsSubscribed)

		case "etc":
			utils.StartSubscribe("etc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Etc.IsSubscribed)

		case "btc":
			utils.StartSubscribe("btc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Btc.IsSubscribed)

		case "ltc":
			utils.StartSubscribe("ltc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Ltc.IsSubscribed)

		case "bch":
			utils.StartSubscribe("bch", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Bch.IsSubscribed)

		case "xlm":
			utils.StartSubscribe("xlm", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Xlm.IsSubscribed)

		case "all":
			utils.StartSubscribe("eth", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Eth.IsSubscribed)
			utils.StartSubscribe("etc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Etc.IsSubscribed)
			utils.StartSubscribe("btc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Btc.IsSubscribed)
			utils.StartSubscribe("ltc", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Ltc.IsSubscribed)
			utils.StartSubscribe("bch", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Bch.IsSubscribed)
			utils.StartSubscribe("xlm", *configData, b, m, &Subscription, Subscription.Info[m.Sender.ID].Xlm.IsSubscribed)

		default:
			b.Send(m.Sender, "Mistake in command!")
		}

	})

	b.Handle("/stop", func(m *tb.Message) {
		params := strings.Split(m.Text, " ")
		switch params[1] {

		case "eth":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Eth.SubsChan, &Subscription, "eth", m, b)

		case "etc":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Etc.SubsChan, &Subscription, "etc", m, b)

		case "btc":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Btc.SubsChan, &Subscription, "btc", m, b)

		case "ltc":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Ltc.SubsChan, &Subscription, "ltc", m, b)

		case "bch":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Bch.SubsChan, &Subscription, "bch", m, b)

		case "xlm":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Xlm.SubsChan, &Subscription, "xlm", m, b)

		case "all":
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Eth.SubsChan, &Subscription, "eth", m, b)
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Etc.SubsChan, &Subscription, "etc", m, b)
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Btc.SubsChan, &Subscription, "btc", m, b)
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Ltc.SubsChan, &Subscription, "ltc", m, b)
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Bch.SubsChan, &Subscription, "bch", m, b)
			utils.StopSubscribe(Subscription.Info[m.Sender.ID].Xlm.SubsChan, &Subscription, "xlm", m, b)

		default:
			b.Send(m.Sender, "Mistake in command!")
		}

	})

	b.Start()
}
