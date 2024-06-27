package structs

// REQUESTS
type RegisterRequest struct {
	EncryptionKey     string `json:"encryption_key" binding:"required"`
	Email             string `json:"email" binding:"required"`
	ReminderTime      string `json:"reminder_time" binding:"required"`
	Timezone          string `json:"timezone" binding:"required"`
	TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
	TelegramUserID    string `json:"telegram_user_id" binding:"required"`
}

type LoginRequest struct {
	Email         string `json:"email" binding:"required"`
	EncryptionKey string `json:"encryption_key" binding:"required"`
}

type ModifyUserRequest struct {
	EncryptionKey     string     `json:"encryption_key" binding:"required"`
	Email             string     `json:"email"`
	ReminderTime      string     `json:"reminder_time" binding:"required"`
	Timezone          string     `json:"timezone" binding:"required"`
	TelegramBotAPIKey string     `json:"telegram_bot_api_key" binding:"required"`
	TelegramUserID    string     `json:"telegram_user_id" binding:"required"`
	Birthdays         []Birthday `json:"birthdays"`
}

type Birthday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// RESPONSES
type Error struct {
	Error string `json:"error"`
}

type Success struct {
	Success bool `json:"success"`
}

type LoginSuccess struct {
	TelegramBotAPIKey string     `json:"telegram_bot_api_key"`
	TelegramUserID    string     `json:"telegram_user_id"`
	ReminderTime      string     `json:"reminder_time"`
	Timezone          string     `json:"timezone"`
	Birthdays         []Birthday `json:"birthdays"`
}

type EncryptionKey struct {
	EncryptionKey string `json:"encryption_key"`
}
