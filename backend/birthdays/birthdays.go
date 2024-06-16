package birthdays

import (
	"fmt"
	"log"
	"time"

	"hbd/encryption"
	"hbd/env"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CallReminderChecker(c *gin.Context) {
	CheckReminders()
}

// checkReminders runs periodically to check for user reminders.
func CheckReminders() {
	println("Checking reminders")
	// now := time.Now().UTC()
	// query := `
	//     SELECT id, telegram_bot_api_key, telegram_user_id FROM users
	//     WHERE reminder_time = $1
	// `

	// rows, err := env.DB.Query(query, now.Format("15:04"))
	rows, err := env.DB.Query("select id, telegram_bot_api_key, telegram_user_id from users")
	if err != nil {
		log.Println("Error querying users:", err)
		return
	}
	defer rows.Close()

	var userId int
	var encryptedBotAPIKey, encryptedUserID string
	for rows.Next() {
		if err := rows.Scan(&userId, &encryptedBotAPIKey, &encryptedUserID); err != nil {
			log.Println("Error scanning user id:", err)
			continue
		}
		botAPIKey, err := encryption.Decrypt(env.MK, encryptedBotAPIKey)
		if err != nil {
			log.Println("Error decrypting bot API key:", err)
			continue
		}
		userID, err := encryption.Decrypt(env.MK, encryptedUserID)
		if err != nil {
			log.Println("Error decrypting user ID:", err)
			continue
		}
		sendBirthdayReminder(userId, botAPIKey, userID)
	}
}

// sendBirthdayReminder sends birthday reminders to the user via Telegram.
func sendBirthdayReminder(userId int, botAPIKey, telegramUserID string) {
	query := `
        SELECT name, date_of_birth FROM birthdays 
        WHERE user_id = $1 AND EXTRACT(MONTH FROM date_of_birth) = $2 AND EXTRACT(DAY FROM date_of_birth) = $3
    `

	now := time.Now().UTC()
	rows, err := env.DB.Query(query, userId, now.Month(), now.Day())
	if err != nil {
		log.Println("Error querying birthdays:", err)
		return
	}
	defer rows.Close()

	var birthdays []string
	for rows.Next() {
		var name string
		var dateOfBirth time.Time
		if err := rows.Scan(&name, &dateOfBirth); err != nil {
			log.Println("Error scanning birthday:", err)
			continue
		}

		age := now.Year() - dateOfBirth.Year()
		if now.Month() < dateOfBirth.Month() || (now.Month() == dateOfBirth.Month() && now.Day() < dateOfBirth.Day()) {
			age--
		}

		birthdays = append(birthdays, fmt.Sprintf("> %s - Turns %d", name, age))
	}

	if len(birthdays) > 0 {
		reminder := fmt.Sprintf("🎂 Birthdays for today: %s\n%s", now.Format("2006-01-02"), formatBirthdays(birthdays))
		log.Println(reminder)
		sendTelegramMessage(botAPIKey, telegramUserID, reminder)
	}
}

// formatBirthdays formats the birthday list into a single string.
func formatBirthdays(birthdays []string) string {
	return fmt.Sprintf("%s", birthdays)
}

// sendTelegramMessage sends a message via the Telegram bot API.
func sendTelegramMessage(botAPIKey, telegramUserID, message string) {
	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Println("Error creating Telegram bot:", err)
		return
	}

	msg := tgbotapi.NewMessageToChannel(telegramUserID, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending Telegram message:", err)
	}
}
