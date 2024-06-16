package structs

type RegisterRequest struct {
	EncryptionKey     string `json:"encryption_key" binding:"required"`
	Email             string `json:"email" binding:"required"`
	ReminderTime      string `json:"reminder_time" binding:"required"`
	Timezone          string `json:"timezone" binding:"required"`
	TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
	TelegramUserID    string `json:"telegram_user_id" binding:"required"`
}

type LoginRequest struct {
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
