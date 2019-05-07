package main

import (
	"log"
	"os"
	"time"

	"encoding/json"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/keyboard"
	"github.com/prazd/nodes_mon_bot/subscribe"
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
	Subscribtion := subscribe.SubNew()

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

	b.Handle("/subscribe", func(m *tb.Message) {
		params := strings.Split(m.Text, " ")

		currencies := []string{"eth", "etc", "xlm", "bch", "btc", "ltc"}

		ok := utils.Contains(params[1], currencies)

		if !ok {
			b.Send(m.Sender, "Mistake in command!")
			return
		}

		utils.StartSubscribe(params[1], *configData, b, m, &Subscribtion)
	})

	b.Handle("/stop", func(m *tb.Message) {
		params := strings.Split(m.Text, " ")
		switch params[1] {
		case "eth":
			close(Subscribtion.Info[m.Sender.ID].Eth.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "eth")
			b.Send(m.Sender, "ETH subscribe stop successful!")
		case "etc":
			close(Subscribtion.Info[m.Sender.ID].Etc.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "etc")
			b.Send(m.Sender, "ETC subscribe stop successful!")
		case "btc":
			close(Subscribtion.Info[m.Sender.ID].Btc.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "btc")
			b.Send(m.Sender, "BTC subscribe stop successful!")
		case "ltc":
			close(Subscribtion.Info[m.Sender.ID].Ltc.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "ltc")
			b.Send(m.Sender, "LTC subscribe stop successful!")
		case "bch":
			close(Subscribtion.Info[m.Sender.ID].Bch.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "bch")
			b.Send(m.Sender, "BCH subscribe stop successful!")
		case "xlm":
			close(Subscribtion.Info[m.Sender.ID].Xlm.SubsChan)
			Subscribtion.Remove(m.Sender.ID, "xlm")
			b.Send(m.Sender, "XLM subscribe stop successful!")
		default:
			b.Send(m.Sender, "Mistake in command!")
		}

	})

	b.Start()
}
