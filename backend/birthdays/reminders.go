package birthdays

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hbd/encryption"
	"hbd/env"
	"hbd/helper"
	"hbd/telegram"
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

	// Define a struct to hold the query result
	type User struct {
		ID                int
		TelegramBotAPIKey string
		TelegramUserID    string
	}

	// Iterate through the rows and scan the data
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.TelegramBotAPIKey, &user.TelegramUserID)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		// Print or log the user struct to see the contents
		fmt.Printf("User: %+v\n", user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	// If you still need to return rows for other uses
	// re-execute the query because rows have already been iterated over
	rows, err = env.DB.Query(query, time.Format("15:04"))
	if err != nil {
		log.Println("Error querying users again:", err)
		return nil, err
	}

	return rows, nil
}

// sendBirthdayReminder sends birthday reminders to the user via Telegram.
func sendBirthdayReminder(userId int, botAPIKey, telegramUserID string) {
	// SQL query to fetch names and dates of birthdays for the given user on the current date
	var query string
	if env.DBType() == "sqlite" {
		query = `
		SELECT name, date FROM birthdays
		WHERE user_id = ? AND 
		cast(strftime('%m', date) as integer) = ? AND 
		cast(strftime('%d', date) as integer) = ?`
	} else {
		query = `
        SELECT name, date FROM birthdays 
        WHERE user_id = $1 AND 
		EXTRACT(MONTH FROM TO_DATE(date, 'YYYY-MM-DD'))::int = $2 AND 
		EXTRACT(DAY FROM TO_DATE(date, 'YYYY-MM-DD'))::int = $3`
	}

	// Get the current date in UTC
	now := time.Now().UTC()
	// Execute the SQL query with user ID and the current month and day as parameters
	rows, err := env.DB.Query(query, userId, now.Month(), now.Day())
	if err != nil {
		log.Println("Error querying birthdays:", err)
		return
	}
	defer rows.Close() // Ensure the rows are closed after processing

	var birthdays []string
	// Iterate over the rows returned by the query
	for rows.Next() {
		var name string
		var date time.Time
		// Scan the name and date fields from the current row
		if err := rows.Scan(&name, &date); err != nil {
			log.Println("Error scanning birthday:", err)
			continue
		}

		// Calculate age only if the year is not 0000
		var ageStr string
		if date.Year() != 0 && date.Year() != now.Year() {
			age := now.Year() - date.Year()
			// Adjust age if the birthday hasn't occurred yet this year
			if now.Month() < date.Month() || (now.Month() == date.Month() && now.Day() < date.Day()) {
				age--
			}
			ageStr = fmt.Sprintf(" - Turns %d", age)
		}

		// Add the birthday info to the list
		birthdays = append(birthdays, fmt.Sprintf("> %s%s", name, ageStr))
	}

	// If there are any birthdays for today, create the reminder message
	if len(birthdays) > 0 {
		reminder := fmt.Sprintf("🎂 Birthdays for today: %s\n\n%s", now.Format("2006-01-02"), helper.JoinStrings(birthdays, "\n"))
		// Send the reminder via Telegram
		telegram.SendTelegramMessage(botAPIKey, telegramUserID, reminder)
	}
}
