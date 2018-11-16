package main

import (
	"os"
	"strconv"
	"time"

	portscanner "github.com/anvie/port-scanner"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("BotToken"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}

	volQ := tb.InlineButton{
		Unique: "ETH",
		Text:   "ETH",
	}
	// Inline
	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{volQ},
	}
	b.Handle("/start", func(m *tb.Message) {
		whitelist := [2]string{os.Getenv("DevOne"), os.Getenv("DevTwo")} // id's
		if strconv.Itoa(m.Sender.ID) == whitelist[0] || strconv.Itoa(m.Sender.ID) == whitelist[1] {
			b.Send(m.Sender, "Привет!Я помогу в мониторинге", &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})

			b.Handle(&volQ, func(c *tb.Callback) {
				ps := portscanner.NewPortScanner("ip", 2*time.Second, 5)
				ethCheck := ps.IsOpen(8545)
				var resp string
				if ethCheck == true {
					resp = "✔"
				} else {
					resp = "✖"
				}
				b.Edit(c.Message, "ETH Node live: "+resp, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})
				b.Respond(c, &tb.CallbackResponse{})
			})
		}
	})
	b.Start()
}
