package structs

import "github.com/dgrijalva/jwt-go"

// REQUESTS
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type RegisterRequest struct {
	Email             string `json:"email" binding:"required" example:"example@lotiguere.com"`
	Password          string `json:"password" binding:"required" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
	ReminderTime      string `json:"reminder_time" binding:"required" example:"15:04"`
	Timezone          string `json:"timezone" binding:"required" example:"America/New_York"`
	TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	TelegramUserID    string `json:"telegram_user_id" binding:"required" example:"123456789"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"example@lotiguere.com"`
	Password string `json:"password" binding:"required" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
}

type ModifyUserRequest struct {
	NewEmail             string `json:"new_email" example:"example2@lotiguere.com"`
	NewPassword          string `json:"new_password" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
	NewReminderTime      string `json:"new_reminder_time" binding:"required" example:"15:04"`
	NewTimezone          string `json:"new_timezone" binding:"required" example:"America/New_York"`
	NewTelegramBotAPIKey string `json:"new_telegram_bot_api_key" binding:"required" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	NewTelegramUserID    string `json:"new_telegram_user_id" binding:"required" example:"123456789"`
}

type BirthdayNameDateModify struct {
	ID   int64  `json:"id" binding:"required" example:"1"`
	Name string `json:"name" binding:"required" example:"John Doe"`
	Date string `json:"date" binding:"required" example:"2021-01-01"`
}

type BirthdayNameDateAdd struct {
	Name string `json:"name" binding:"required" example:"John Doe"`
	Date string `json:"date" binding:"required" example:"2021-01-01"`
}

type BirthdayFull struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"John Doe"`
	Date string `json:"date" example:"2021-01-01"`
}

type BirthdayID struct {
	ID int64 `json:"id" example:"1"`
}

// RESPONSES
type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success bool `json:"success"`
}

type LoginSuccess struct {
	Token             string         `json:"token"`
	TelegramBotAPIKey string         `json:"telegram_bot_api_key" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	TelegramUserID    string         `json:"telegram_user_id" example:"123456789"`
	ReminderTime      string         `json:"reminder_time" example:"15:04"`
	Timezone          string         `json:"timezone" example:"America/New_York"`
	Birthdays         []BirthdayFull `json:"birthdays"`
}

type UserData struct {
	ID                int64          `json:"id" example:"1"`
	TelegramBotAPIKey string         `json:"telegram_bot_api_key" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	TelegramUserID    string         `json:"telegram_user_id" example:"123456789"`
	ReminderTime      string         `json:"reminder_time" example:"15:04"`
	Timezone          string         `json:"timezone" example:"America/New_York"`
	Birthdays         []BirthdayFull `json:"birthdays"`
}

type Password struct {
	Password string `json:"password" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
}

// User data for when a query is made for reminders
// UserData holds the decrypted user data
type UserDataShort struct {
	UserID         int
	BotAPIKey      string
	TelegramUserID string
}
