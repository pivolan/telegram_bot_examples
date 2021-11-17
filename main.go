package main

import (
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

const TELEGRAM_BOT_API = "2124047747:AAFjAHlV37rFCbjqZ9c6CmqtCREUTYOydQE"

var botMessageId int
var userMessageId int

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
			if update.InlineQuery != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "inline query: "+update.InlineQuery.Query)
				bot.Send(msg)
				return
			}
			if update.CallbackQuery != nil { //on button clicked
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "calback query: "+update.CallbackQuery.Data)
				bot.Send(msg)
				return
			}

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
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(msg)
					return
				case "buttons": //show buttons inline keyboard
					kb := InlineKeyboard()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "inline keyboard")
					msg.ReplyMarkup = kb
					bot.Send(msg)
					return
				case "markdown": //markdown message smiles
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("*Нажмите* **curs** _italy_\n\n# header\n\n 🙂 some text \xF0\x9F\x98\x81 some text"))
					msg.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(msg)
					return
				case "file": //edit bot message
					data := tgbotapi.FileBytes{Name: "filename.txt", Bytes: []byte("empty file")}
					msg := tgbotapi.NewDocumentUpload(message.Chat.ID, data)
					msg.Caption = "caption"
					bot.Send(msg)
					return
				case "edit": //edit bot message
					msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, botMessageId, fmt.Sprintf("message changed"))
					msg.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(msg)
					return
				case "edit_self": //user message cannot be edited
					msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, userMessageId, fmt.Sprintf("your message changed"))
					msg.ParseMode = tgbotapi.ModeMarkdown
					_, err := bot.Send(msg)
					if err != nil {
						fmt.Println(err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("error edit message: %s", err))
						bot.Send(msg)
					}
					return
				default:
					commands := []string{"/k", "/khide", "/buttons", "/markdown", "/edit", "/edit_self", "/file"}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("command getted: %s\n commands allowed:\n%s", update.Message.Command(), strings.Join(commands, "\n")))
					bot.Send(msg)
				}
				return
			}
			switch message.Text {
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Text: %s, ChatID: %d, username: %s", message.Text, message.Chat.ID, message.Chat.UserName))
				r, err := bot.Send(msg)
				if err != nil {
					return
				}
				botMessageId = r.MessageID
				userMessageId = message.MessageID
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
