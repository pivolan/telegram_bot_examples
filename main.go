package main

import (
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

const TELEGRAM_BOT_API = "2124047747:AAFjAHlV37rFCbjqZ9c6CmqtCREUTYOydQE"

func main() {
	bot, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_API)
	if err != nil {
		log.Fatalf("can not init telegram bot, err: %s", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update) {
			message := update.Message

			if message.IsCommand() {
				args := strings.TrimSpace(update.Message.CommandArguments())
				parts := strings.Split(args, " ")
				if len(parts) == 1 && parts[0] == "" {
					parts = []string{}
				}
				switch update.Message.Command() {
				case "k": //show keyboard
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Клава"))
					msg.ReplyMarkup = DefaultKeyboard()
					bot.Send(msg)
					return
				case "khide": //hide keyboard
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Скрыл клаву"))
					msg.ReplyMarkup = DefaultKeyboard()
					bot.Send(msg)
					return
				case "buttons": //show buttons inline keyboard
					kb := InlineKeyboard()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "inline keyboard")
					msg.ReplyMarkup = kb
					bot.Send(msg)
					return
				case "markdown": //markdown message smiles
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("*Нажмите* **curs** _italy_\n# header 🙂 some text \xF0\x9F\x98\x81 some text"))
					msg.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(msg)
					return
				default:
				}
				return
			}

			switch message.Text {
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Нажмите Next чтобы получить следующую историю"))
				bot.Send(msg)
				return
			}
		}(update)
	}

}

func InlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := map[string]string{"\xF0\x9F\x98\x81 some text": "/some_route", "12": "some text"}

	keyboardButtons := []tgbotapi.InlineKeyboardButton{}
	for buttonName, buttonValue := range buttons {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(buttonName, buttonValue))
	}
	kb := tgbotapi.NewInlineKeyboardMarkup(keyboardButtons)

	return kb
}
func DefaultKeyboard() tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{}
	buttons := [][]string{
		{"11", "12", "13"},
		{"21", "22", "23"},
	}
	for _, cells := range buttons {
		row := []tgbotapi.KeyboardButton{}
		for _, cell := range cells {
			row = append(row, tgbotapi.KeyboardButton{cell, false, false})
		}
		rows = append(rows, row)
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}
