package keyboard

import tb "gopkg.in/tucnak/telebot.v2"

var (
	EthButton = tb.ReplyButton{Text: "ETH"}

	EtcButton = tb.ReplyButton{Text: "ETC"}

	BtcButton = tb.ReplyButton{Text: "BTC"}

	BchButton = tb.ReplyButton{Text: "BCH"}

	LtcButton = tb.ReplyButton{Text: "LTC"}

	SubscriptionStatus = tb.ReplyButton{Text: "Subscription status"}

	MainMenu = [][]tb.ReplyButton{
		[]tb.ReplyButton{SubscriptionStatus},
		[]tb.ReplyButton{EthButton, EtcButton, BtcButton},
		[]tb.ReplyButton{BchButton, LtcButton},
	}
)
