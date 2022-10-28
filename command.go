package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const (
	bashCommand   = "bash"
	vpnCommand    = "pivpn"
	uptimeCommand = "uptime"
	rebootCommand = "reboot"

	cpuBotCommand  = "cpu"
	memBotCommand  = "mem"
	diskBotCommand = "disk"
	vpnBotCommand  = "vpn"
	defaultAnswer  = "Oops!"
)

var (
	cpuCommand  = []string{"-c", `top -bn2 | grep '%Cpu' | tail -1 | grep -P '(....|...) id,'| awk '{print "CPU usage: " 100-$8 "%"}'`}
	memCommand  = []string{"-c", `top -bn2 | grep 'MiB Mem' | tail -1 | awk '{printf "Memory usage: %s total; %s free; %s used\n",$4,$6,$8}'`}
	diskCommand = []string{"df", "-h"}
)

var (
	removeVPNOutput = []string{
		"[1m",
		"[0m",
		"[4mName",
		"[4mLast Seen",
		"[4mRemote IP",
		"[4mBytes Sent",
		"[4mVirtual IP",
		"[4mBytes Received",
		"::: Connected Clients List :::",
		"::: Disabled clients :::",
	}
)

func run(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, userID int64) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			log.WithFields(log.Fields{
				"id":         update.SentFrom().ID,
				"first name": update.SentFrom().FirstName,
				"user name":  update.SentFrom().UserName,
			}).Info("not command")

			continue
		}

		if update.SentFrom().ID != userID {
			log.WithFields(log.Fields{
				"id":         update.SentFrom().ID,
				"first name": update.SentFrom().FirstName,
				"user name":  update.SentFrom().UserName,
			}).Info("wrong user")

			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case cpuBotCommand:
			cmd := exec.Command(bashCommand, cpuCommand...)
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.WithError(err).Fatal("run cpu")

				msg.Text = defaultAnswer
			} else {
				msg.Text = outb.String()
			}
		case memBotCommand:
			cmd := exec.Command(bashCommand, memCommand...)
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.WithError(err).Fatal("run mem")

				msg.Text = defaultAnswer
			} else {
				msg.Text = outb.String()
			}
		case diskBotCommand:
			cmd := exec.Command(diskCommand[0], diskCommand[1])
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.WithError(err).Fatal("run disk")

				msg.Text = defaultAnswer
			} else {
				msg.Text = outb.String()
			}
		case uptimeCommand:
			cmd := exec.Command(uptimeCommand)
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.WithError(err).Fatal("run uptime")

				msg.Text = defaultAnswer
			} else {
				msg.Text = outb.String()
			}
		case vpnBotCommand:
			cmd := exec.Command(vpnCommand, "-c")
			var outb bytes.Buffer
			cmd.Stdout = &outb

			if err := cmd.Run(); err != nil {
				log.WithError(err).Fatal("run vpn")

				msg.Text = defaultAnswer
			} else {
				res := outb.String()

				for _, str := range removeVPNOutput {
					res = strings.Replace(res, str, "", -1)
				}

				msg.ParseMode = tgbotapi.ModeMarkdownV2
				msg.Text = fmt.Sprintf("```%v```", res)
			}
		case rebootCommand:
			exec.Command(rebootCommand).Run()
		default:
			msg.Text = defaultAnswer
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
