package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendTelegramMessage sends a message via the Telegram bot API.
func SendTelegramMessage(botAPIKey, telegramUserID, message string) {
	// Create a new Telegram bot instance using the provided API key
	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Println("Error creating Telegram bot:", err)
		return
	}

	// Create a new message to send to the specified Telegram user/channel
	msg := tgbotapi.NewMessageToChannel(telegramUserID, message)
	// Send the message via the Telegram bot API
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending Telegram message:", err)
	}
}
