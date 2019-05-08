package keyboard

import tb "gopkg.in/tucnak/telebot.v2"

var (
	EthButton = tb.ReplyButton{Text: "ETH"}

	EtcButton = tb.ReplyButton{Text: "ETC"}

	BtcButton = tb.ReplyButton{Text: "BTC"}

	BchButton = tb.ReplyButton{Text: "BCH"}

	LtcButton = tb.ReplyButton{Text: "LTC"}

	XlmButton = tb.ReplyButton{Text: "XLM"}

	SubscribeStatus = tb.ReplyButton{Text: "Subscribe Status"}

	MainMenu = [][]tb.ReplyButton{
		[]tb.ReplyButton{SubscribeStatus},
		[]tb.ReplyButton{EthButton, EtcButton, BtcButton},
		[]tb.ReplyButton{XlmButton, BchButton, LtcButton},
	}
)
