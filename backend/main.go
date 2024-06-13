package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/chacha20poly1305"

	"hbd/models"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=postgres port=6684 host=localhost sslmode=disable")
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

func generateEncryptionKey() (string, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func getEncryptionKey(c *gin.Context) {
	key, err := generateEncryptionKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate encryption key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"encryption_key": key})
}

func register(c *gin.Context) {
	type RegisterRequest struct {
		EncryptionKey     string `json:"encryption_key" binding:"required"`
		ReminderTime      string `json:"reminder_time" binding:"required"` // Expecting time in the format "15:04" (HH:MM)
		Timezone          string `json:"timezone" binding:"required"`      // Expecting IANA Timezone format like "America/New_York"
		TelegramBotAPIKey string `json:"telegram_bot_api_key" binding:"required"`
		TelegramUserID    string `json:"telegram_user_id" binding:"required"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedBotAPIKey, err := encrypt(req.EncryptionKey, req.TelegramBotAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt Telegram bot API key"})
		return
	}

	encryptedUserID, err := encrypt(req.EncryptionKey, req.TelegramUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt Telegram user ID"})
		return
	}

	// Parse reminder time in user's local timezone
	location, err := time.LoadLocation(req.Timezone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timezone"})
		return
	}

	// Assuming reminder time is for the current date in the provided timezone
	now := time.Now()
	reminderTime, err := time.ParseInLocation("15:04", req.ReminderTime, location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder time format"})
		return
	}

	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).UTC()

	user := models.User{
		ID:                uuid.New().String(),
		ReminderTime:      reminderTime,
		TelegramBotAPIKey: hex.EncodeToString(encryptedBotAPIKey),
		TelegramUserID:    hex.EncodeToString(encryptedUserID),
	}

	err = user.Insert(c, db, boil.Infer())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.FindUser(c, db, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	decryptedBotAPIKey, err := decrypt(req.EncryptionKey, user.TelegramBotAPIKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid encryption key"})
		return
	}

	decryptedUserID, err := decrypt(req.EncryptionKey, user.TelegramUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid encryption key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"telegram_bot_api_key": decryptedBotAPIKey, "telegram_user_id": decryptedUserID})
}

func updateBirthday(c *gin.Context) {
	// Implementation for updating birthday details
}

func encrypt(encryptionKey string, plaintext string) ([]byte, error) {
	key, err := hex.DecodeString(encryptionKey)
	if err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plaintext)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

func decrypt(encryptionKey string, ciphertextHex string) (string, error) {
	key, err := hex.DecodeString(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
