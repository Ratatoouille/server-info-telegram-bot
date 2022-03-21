package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	bashCommand = "bash"
)

var (
	cpuCommand  = []string{"-c", `top -bn2 | grep '%Cpu' | tail -1 | grep -P '(....|...) id,'| awk '{print "CPU usage: " 100-$8 "%"}'`}
	memCommand  = []string{"-c", `top -bn2 | grep 'MiB Mem' | tail -1 | awk '{printf "Memory usage: %s total; %s free; %s used\n",$4,$6,$8}'`}
	diskCommand = []string{"df", "-h"}
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	userID, err := strconv.ParseInt(os.Getenv("USER_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		if update.SentFrom().ID != userID {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "cpu":
			cmd := exec.Command(bashCommand, cpuCommand...)
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.Println(err)
				msg.Text = "Oops!"
			} else {
				msg.Text = outb.String()
			}
		case "mem":
			cmd := exec.Command(bashCommand, memCommand...)
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.Println(err)
				msg.Text = "Oops!"
			} else {
				msg.Text = outb.String()
			}
		case "disk":
			cmd := exec.Command(diskCommand[0], diskCommand[1])
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.Println(err)
				msg.Text = "Oops!"
			} else {
				msg.Text = outb.String()
			}
		default:
			msg.Text = "Oops!"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
