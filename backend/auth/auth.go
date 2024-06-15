package auth

import (
	"crypto/rand"
	"encoding/hex"
	"hbd/encryption"
	"hbd/env"
	"hbd/helper"
	"hbd/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetEncryptionKey(c *gin.Context) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to generate encryption key", false) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"encryption_key": hex.EncodeToString(key)})
}

func Register(c *gin.Context) {
	type RegisterRequest struct {
		EncryptionKey     string `json:"encryption_key" binding:"required"`
		Email             string `json:"email" binding:"required"`
		ReminderTime      string `json:"reminder_time" binding:"required"`
		Timezone          string `json:"timezone" binding:"required"`
		TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
		TelegramUserID    string `json:"telegram_user_id" binding:"required"`
	}

	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return
	}

	// Hash the email to check for uniqueness
	emailHash := encryption.HashStringWithSHA256(req.Email)

	// Check if the email hash already exists in the database
	exists, err := models.Users(models.UserWhere.EmailHash.EQ(emailHash)).Exists(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Bot API key hash
	botAPIKeyHash := encryption.HashStringWithSHA256(req.TelegramBotAPIKey)

	// Check if the bot API key hash already exists in the database
	exists, err = models.Users(models.UserWhere.TelegramBotAPIKeyHash.EQ(botAPIKeyHash)).Exists(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Telegram bot API key already registered"})
		return
	}

	encryptedBotAPIKey, err := encryption.Encrypt(env.MK, req.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram bot API key", false) {
		return
	}

	// User ID Hash
	telegramUserIDHash := encryption.HashStringWithSHA256(req.TelegramUserID)

	// Check if the Telegram user ID hash already exists in the database
	exists, err = models.Users(models.UserWhere.TelegramUserIDHash.EQ(telegramUserIDHash)).Exists(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Telegram user ID already registered"})
		return
	}

	encryptedUserID, err := encryption.Encrypt(env.MK, req.TelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram user ID", false) {
		return
	}

	location, err := time.LoadLocation(req.Timezone)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid timezone", false) {
		return
	}

	now := time.Now()
	reminderTime, err := time.ParseInLocation("15:04", req.ReminderTime, location)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid reminder time format", false) {
		return
	}

	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).UTC()

	// Encrypt the user's key using the master key before storing it
	encryptedKey, err := encryption.Encrypt(env.MK, req.EncryptionKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt encryption key", false) {
		return
	}

	user := models.User{
		EmailHash:             emailHash,
		EncryptionKey:         hex.EncodeToString(encryptedKey),
		EncryptionKeyHash:     encryption.HashStringWithSHA256(req.EncryptionKey),
		ReminderTime:          reminderTime,
		Timezone:              req.Timezone,
		TelegramBotAPIKey:     hex.EncodeToString(encryptedBotAPIKey),
		TelegramBotAPIKeyHash: botAPIKeyHash,
		TelegramUserID:        hex.EncodeToString(encryptedUserID),
		TelegramUserIDHash:    telegramUserIDHash,
	}

	err = user.Insert(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to create user", false) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Login(c *gin.Context) {
	type LoginRequest struct {
		EncryptionKey string `json:"encryption_key" binding:"required"`
	}
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return
	}

	// Fetch the user with the given encryption key
	user, err := models.Users(models.UserWhere.EncryptionKeyHash.EQ(encryption.HashStringWithSHA256(req.EncryptionKey))).One(c, env.DB)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	decryptedBotAPIKey, err := encryption.Decrypt(env.MK, user.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	decryptedUserID, err := encryption.Decrypt(env.MK, user.TelegramUserID)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	// show reminder time in the user's designated timezone
	reminderTime := user.ReminderTime.In(time.FixedZone(user.Timezone, 0))

	// send the entire birthday list
	// find the birthdays by user id
	birthdays, err := models.Birthdays(models.BirthdayWhere.UserID.EQ(user.ID)).All(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to fetch birthdays", false) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"telegram_bot_api_key": decryptedBotAPIKey,
		"telegram_user_id":     decryptedUserID,
		"reminder_time":        reminderTime,
		"timezone":             user.Timezone,
		"birthdays":            birthdays,
	})
}

func DeleteUser(c *gin.Context) {
	type DeleteUserRequest struct {
		EncryptionKey string `json:"encryption_key" binding:"required"`
	}
	var req DeleteUserRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return
	}

	// Fetch the user with the given encryption key
	user, err := models.Users(models.UserWhere.EncryptionKeyHash.EQ(encryption.HashStringWithSHA256(req.EncryptionKey))).One(c, env.DB)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	_, err = user.Delete(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to delete user", false) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ModifyUser(c *gin.Context) {
	type ModifyUserRequest struct {
		EncryptionKey     string `json:"encryption_key" binding:"required"`
		Email             string `json:"email"`
		ReminderTime      string `json:"reminder_time" binding:"required"`
		Timezone          string `json:"timezone" binding:"required"`
		TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
		TelegramUserID    string `json:"telegram_user_id" binding:"required"`
	}
	var req ModifyUserRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return
	}

	// Fetch the user with the given encryption key
	println(encryption.HashStringWithSHA256(req.EncryptionKey))
	user, err := models.Users(models.UserWhere.EncryptionKeyHash.EQ(encryption.HashStringWithSHA256(req.EncryptionKey))).One(c, env.DB)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	location, err := time.LoadLocation(req.Timezone)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid timezone", false) {
		return
	}

	now := time.Now()
	reminderTime, err := time.ParseInLocation("15:04", req.ReminderTime, location)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid reminder time format", false) {
		return
	}

	telegramBotAPIKeyHash := encryption.HashStringWithSHA256(req.TelegramBotAPIKey)
	encryptedBotAPIKey, err := encryption.Encrypt(env.MK, req.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram bot API key", false) {
		return
	}

	telegramUserIDHash := encryption.HashStringWithSHA256(req.TelegramUserID)
	encryptedUserID, err := encryption.Encrypt(env.MK, req.TelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram user ID", false) {
		return
	}

	// hash the user's email (to be updated)
	emailHash := encryption.HashStringWithSHA256(req.Email)

	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).UTC()

	user.EmailHash = emailHash
	user.ReminderTime = reminderTime
	user.Timezone = req.Timezone
	user.TelegramBotAPIKey = hex.EncodeToString(encryptedBotAPIKey)
	user.TelegramBotAPIKeyHash = telegramBotAPIKeyHash
	user.TelegramUserID = hex.EncodeToString(encryptedUserID)
	user.TelegramUserIDHash = telegramUserIDHash

	_, err = user.Update(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to update user", false) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
