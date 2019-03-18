package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"encoding/json"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/handlers"
	"github.com/prazd/nodes_mon_bot/keyboard"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"path/filepath"
)

func ReadConfig() config.Config {

	defaultConfigPath, _ := filepath.Abs("../config/config.json")
	configFile, err := os.Open(defaultConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var conf config.Config

	json.Unmarshal(byteValue, &conf)

	return conf
}

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
				InlineKeyboard: keyboard.MainMenu,
			})

			configData := ReadConfig()

			// Check ETH nodes status
			b.Handle(&keyboard.EthButton, func(c *tb.Callback) {

				message := handlers.IsAlive("eth", configData)

				b.Edit(c.Message, message, &tb.ReplyMarkup{
					InlineKeyboard: keyboard.MainMenu,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			// Check ETC nodes status
			b.Handle(&keyboard.EtcButton, func(c *tb.Callback) {

				message := handlers.IsAlive("etc", configData)

				b.Edit(c.Message, message, &tb.ReplyMarkup{
					InlineKeyboard: keyboard.MainMenu,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			// Check BTC nodes status
			b.Handle(&keyboard.BtcButton, func(c *tb.Callback) {
				message := handlers.IsAlive("btc", configData)

				b.Edit(c.Message, message, &tb.ReplyMarkup{
					InlineKeyboard: keyboard.MainMenu,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			// Check BCH nodes status
			b.Handle(&keyboard.BchButton, func(c *tb.Callback) {
				message := handlers.IsAlive("bch", configData)

				b.Edit(c.Message, message, &tb.ReplyMarkup{
					InlineKeyboard: keyboard.MainMenu,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			// Check LTC nodes status
			b.Handle(&keyboard.LtcButton, func(c *tb.Callback) {
				message := handlers.IsAlive("ltc", configData)

				b.Edit(c.Message, message, &tb.ReplyMarkup{
					InlineKeyboard: keyboard.MainMenu,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

		}
	})

	b.Start()
}
