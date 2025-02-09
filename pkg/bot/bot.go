package bot

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func MustNew(token string) *tele.Bot {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("error creating bot: %v", err)
	}

	return bot
}
