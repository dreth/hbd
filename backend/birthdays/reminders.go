package birthdays

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hbd/encryption"
	"hbd/env"
	"hbd/helper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CheckReminders runs periodically to check for user reminders.
func CheckReminders() {
	now := time.Now().UTC()
	rows, err := queryWithTime(now)
	if err != nil {
		log.Println("Error querying users:", err)
		return
	}
	remindBirthdays(rows)
}

// remindBirthdays
func remindBirthdays(rows *sql.Rows) {
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

// Query with time
func queryWithTime(time time.Time) (*sql.Rows, error) {
	log.Printf("Checking reminders. Timestamp: %s", time)
	query := `
	    SELECT id, telegram_bot_api_key, telegram_user_id FROM users
	    WHERE reminder_time = $1
	`
	rows, err := env.DB.Query(query, time.Format("15:04"))
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	return rows, nil
}

// sendBirthdayReminder sends birthday reminders to the user via Telegram.
func sendBirthdayReminder(userId int, botAPIKey, telegramUserID string) {
	query := `
        SELECT name, date FROM birthdays 
        WHERE user_id = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(DAY FROM date) = $3
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
		var date time.Time
		if err := rows.Scan(&name, &date); err != nil {
			log.Println("Error scanning birthday:", err)
			continue
		}

		// Calculate age only if the year is not 0000
		var ageStr string
		if date.Year() != 0 && date.Year() != now.Year() {
			age := now.Year() - date.Year()
			if now.Month() < date.Month() || (now.Month() == date.Month() && now.Day() < date.Day()) {
				age--
			}
			ageStr = fmt.Sprintf(" - Turns %d", age)
		}

		birthdays = append(birthdays, fmt.Sprintf("> %s%s", name, ageStr))
	}

	if len(birthdays) > 0 {
		reminder := fmt.Sprintf("🎂 Birthdays for today: %s\n\n%s", now.Format("2006-01-02"), helper.JoinStrings(birthdays, "\n"))
		sendTelegramMessage(botAPIKey, telegramUserID, reminder)
	}
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
