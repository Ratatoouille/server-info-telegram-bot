package main

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const (
	offset  = 0
	timeout = 60
)

func main() {
	log.New().SetFormatter(
		&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.WithError(err).Fatal("can't create bot api")
	}

	// bot.Debug = true

	u := tgbotapi.NewUpdate(offset)
	u.Timeout = timeout

	updates := bot.GetUpdatesChan(u)

	userID, err := strconv.ParseInt(os.Getenv("USER_ID"), 10, 64)
	if err != nil {
		log.WithError(err).Fatal("can't parse user id")
	}

	run(bot, updates, userID)
}
