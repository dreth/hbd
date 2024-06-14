package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"hbd/encryption"
	"hbd/helper"
	"hbd/models"
)

var db *sql.DB
var masterKey string
var databaseURL string

func main() {
	var err error

	// Load the MASTER_KEY from the environment
	masterKey = os.Getenv("MASTER_KEY")
	if masterKey == "" {
		log.Fatal("MASTER_KEY environment variable not set")
	}

	// Load the DATABASE_URL from the environment
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	boil.SetDB(db)

	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/update-birthday", updateBirthday)
	router.GET("/generate-encryption-key", getEncryptionKey)

	router.Run(":8080")
}

func getEncryptionKey(c *gin.Context) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to generate encryption key") {
		return
	}
	c.JSON(http.StatusOK, gin.H{"encryption_key": hex.EncodeToString(key)})
}

func register(c *gin.Context) {
	type RegisterRequest struct {
		Key               string `json:"key" binding:"required"`
		Email             string `json:"email" binding:"required"`
		ReminderTime      string `json:"reminder_time" binding:"required"`
		Timezone          string `json:"timezone" binding:"required"`
		TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
		TelegramUserID    string `json:"telegram_user_id" binding:"required"`
	}

	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request") {
		return
	}

	encryptedEmail, err := encryption.Encrypt(masterKey, req.Email)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt email") {
		return
	}

	encryptedBotAPIKey, err := encryption.Encrypt(req.Key, req.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram bot API key") {
		return
	}

	encryptedUserID, err := encryption.Encrypt(req.Key, req.TelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram user ID") {
		return
	}

	location, err := time.LoadLocation(req.Timezone)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid timezone") {
		return
	}

	now := time.Now()
	reminderTime, err := time.ParseInLocation("15:04", req.ReminderTime, location)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid reminder time format") {
		return
	}

	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).UTC()

	// Encrypt the user's key using the master key before storing it
	encryptedKey, err := encryption.Encrypt(masterKey, req.Key)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt encryption key") {
		return
	}

	user := models.User{
		ID:                uuid.New().String(),
		Email:             hex.EncodeToString(encryptedEmail),
		EncryptionKey:     hex.EncodeToString(encryptedKey),
		ReminderTime:      reminderTime,
		TelegramBotAPIKey: hex.EncodeToString(encryptedBotAPIKey),
		TelegramUserID:    hex.EncodeToString(encryptedUserID),
	}

	err = user.Insert(c, db, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to create user") {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}

func login(c *gin.Context) {
	type LoginRequest struct {
		UserID        string `json:"user_id" binding:"required"`
		EncryptionKey string `json:"encryption_key" binding:"required"`
	}
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if helper.HE(c, err, http.StatusBadRequest, err.Error()) {
		return
	}

	user, err := models.FindUser(c, db, req.UserID)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid user ID") {
		return
	}

	// Decrypt the stored encrypted encryption key using the master key
	decryptedKey, err := encryption.Decrypt(masterKey, user.EncryptionKey)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid master encryption key") {
		return
	}

	// Ensure the provided encryption key matches the decrypted stored key
	if decryptedKey != req.EncryptionKey {
		helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key")
		return
	}

	decryptedBotAPIKey, err := encryption.Decrypt(req.EncryptionKey, user.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key") {
		return
	}

	decryptedUserID, err := encryption.Decrypt(req.EncryptionKey, user.TelegramUserID)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key") {
		return
	}

	c.JSON(http.StatusOK, gin.H{"telegram_bot_api_key": decryptedBotAPIKey, "telegram_user_id": decryptedUserID})
}

func updateBirthday(c *gin.Context) {
	// Implementation for updating birthday details
}
