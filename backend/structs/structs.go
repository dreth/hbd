package structs

// REQUESTS
type RegisterRequest struct {
	EncryptionKey     string `json:"encryption_key" binding:"required" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
	Email             string `json:"email" binding:"required" example:"example@hbd.wtf"`
	ReminderTime      string `json:"reminder_time" binding:"required" example:"15:04"`
	Timezone          string `json:"timezone" binding:"required" example:"America/New_York"`
	TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	TelegramUserID    string `json:"telegram_user_id" binding:"required" example:"123456789"`
}

type LoginRequest struct {
	Email         string `json:"email" binding:"required" example:"example@hbd.wtf"`
	EncryptionKey string `json:"encryption_key" binding:"required" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
}

type ModifyUserRequest struct {
	Auth                 LoginRequest `json:"auth" binding:"required"`
	NewEmail             string       `json:"new_email" example:"example2@hbd.wtf"`
	NewReminderTime      string       `json:"new_reminder_time" binding:"required" example:"15:04"`
	NewTimezone          string       `json:"new_timezone" binding:"required" example:"America/New_York"`
	NewTelegramBotAPIKey string       `json:"new_telegram_bot_api_key" binding:"required" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	NewTelegramUserID    string       `json:"new_telegram_user_id" binding:"required" example:"123456789"`
}

type BirthdayNameDateModify struct {
	Auth LoginRequest `json:"auth" binding:"required"`
	ID   int          `json:"id" binding:"required" example:"1"`
	Name string       `json:"name" binding:"required" example:"John Doe"`
	Date string       `json:"date" binding:"required" example:"2021-01-01"`
}

type BirthdayNameDateAdd struct {
	Auth LoginRequest `json:"auth" binding:"required"`
	Name string       `json:"name" binding:"required" example:"John Doe"`
	Date string       `json:"date" binding:"required" example:"2021-01-01"`
}

type BirthdayFull struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"John Doe"`
	Date string `json:"date" example:"2021-01-01"`
}

type BirthdayID struct {
	ID int `json:"id" example:"1"`
}

// RESPONSES
type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success bool `json:"success"`
}

type LoginSuccess struct {
	TelegramBotAPIKey string         `json:"telegram_bot_api_key" example:"270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"`
	TelegramUserID    string         `json:"telegram_user_id" example:"123456789"`
	ReminderTime      string         `json:"reminder_time" example:"15:04"`
	Timezone          string         `json:"timezone" example:"America/New_York"`
	Birthdays         []BirthdayFull `json:"birthdays"`
}

type EncryptionKey struct {
	EncryptionKey string `json:"encryption_key" example:"9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"`
}

// User data for when a query is made for reminders
// UserData holds the decrypted user data
type UserData struct {
	UserID         int
	BotAPIKey      string
	TelegramUserID string
}
