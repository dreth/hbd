package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hbd/encryption"
	"hbd/env"
	"hbd/helper"
	"hbd/models"
	"hbd/structs"
	"hbd/telegram"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Authenticate authenticates a user based on the request payload.
// It supports requests with direct authentication fields as well as nested auth objects.
// The function binds the JSON payload to the provided request struct, extracts the email
// and encryption key, hashes them, and fetches the corresponding user from the database.
//
// Parameters:
//   - c: The Gin context, which provides request and response handling.
//   - req: A pointer to the request struct which will be populated with the JSON payload.
//
// Returns:
//   - *models.User: The authenticated user object if authentication is successful.
//   - error: An error object if any error occurs during the process.
func Authenticate(c *gin.Context, req interface{}) (*models.User, error) {
	// Bind JSON request to the provided struct
	err := c.ShouldBindJSON(req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return nil, err
	}

	// Variables to hold email and encryption key
	var email, encryptionKey string

	// Use reflection to get the underlying value of the request struct
	reqValue := reflect.ValueOf(req).Elem()

	// Try to find an 'Auth' field in the request struct
	authField := reqValue.FieldByName("Auth")

	// Check if 'Auth' field is valid and is a struct
	if authField.IsValid() && authField.Kind() == reflect.Struct {
		// Extract email and encryption key from the 'Auth' struct
		email = authField.FieldByName("Email").String()
		encryptionKey = authField.FieldByName("EncryptionKey").String()
	} else {
		// Extract email and encryption key from the top-level struct
		email = reqValue.FieldByName("Email").String()
		encryptionKey = reqValue.FieldByName("EncryptionKey").String()
	}

	// Hash the email and encryption key
	emailHash := encryption.HashStringWithSHA256(email)
	encryptionKeyHash := encryption.HashStringWithSHA256(encryptionKey)

	// Fetch the user with the given encryption key and email hash from the database
	user, err := models.Users(
		qm.Where("email_hash = ?", emailHash),
		qm.Where("encryption_key_hash = ?", encryptionKeyHash),
	).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		return nil, err
	}

	// Return the authenticated user
	return user, nil
}

// @Summary Generate a new encryption key
// @Description This endpoint generates a new encryption key for the user.
// @Produce  json
// @Success 200 {object} structs.EncryptionKey
// @Failure 500 {object} structs.Error "Failed to generate encryption key"
// @Router /generate-encryption-key [get]
// @Tags auth
// @x-order 1
func GetEncryptionKey(c *gin.Context) {
	// Create a byte slice to hold the generated encryption key
	key := make([]byte, 32)

	// Generate a random 32-byte encryption key
	_, err := rand.Read(key)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to generate encryption key", false) {
		return
	}

	// Respond with the generated encryption key in hexadecimal format
	c.JSON(http.StatusOK, structs.EncryptionKey{EncryptionKey: hex.EncodeToString(key)})
}

// @Summary Register a new user
// @Description This endpoint registers a new user with their email, Telegram bot API key, and other details.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.RegisterRequest  true  "Register user"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 409 {object} structs.Error "Email or Telegram bot API key already registered"
// @Failure 500 {object} structs.Error "Failed to create user"
// @Router /register [post]
// @Tags auth
// @x-order 2
func Register(c *gin.Context) {

	var req structs.RegisterRequest
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

	// Hash the encryption key to check for uniqueness
	encryptionKeyHash := encryption.HashStringWithSHA256(req.EncryptionKey)
	if exists, err = models.Users(models.UserWhere.EncryptionKeyHash.EQ(encryptionKeyHash)).Exists(c, env.DB); helper.HE(c, err, http.StatusInternalServerError, "Failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Encryption key already registered"})
		return
	}

	encryptedBotAPIKey, err := encryption.Encrypt(env.MK, req.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram bot API key", false) {
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

	// Ensure reminderTime is in UTC
	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).In(time.UTC)

	// Encrypt the user's key using the master key before storing it
	encryptedKey, err := encryption.Encrypt(env.MK, req.EncryptionKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt encryption key", false) {
		return
	}

	// Create a new user object
	user := models.User{
		EmailHash:             emailHash,
		EncryptionKey:         hex.EncodeToString(encryptedKey),
		EncryptionKeyHash:     encryption.HashStringWithSHA256(req.EncryptionKey),
		ReminderTime:          reminderTime.String(),
		Timezone:              req.Timezone,
		TelegramBotAPIKey:     hex.EncodeToString(encryptedBotAPIKey),
		TelegramBotAPIKeyHash: encryption.HashStringWithSHA256(req.TelegramUserID),
		TelegramUserID:        hex.EncodeToString(encryptedUserID),
		TelegramUserIDHash:    encryption.HashStringWithSHA256(req.TelegramUserID),
	}

	err = user.Insert(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to create user", false) {
		return
	}

	// As the user was successfully created, send a telegram message through the bot and ID to confirm the registration
	telegram.SendTelegramMessage(req.TelegramBotAPIKey, req.TelegramUserID, fmt.Sprintf("🎂 Your user has been successfully registered, through this bot and user ID you'll receive your birthday reminders (if there's any) at %s (Timezone: %s).", req.ReminderTime, req.Timezone))

	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// @Summary Login a user
// @Description This endpoint logs in a user by validating their email and encryption key.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.LoginRequest  true  "Login user"
// @Success 200 {object} structs.LoginSuccess
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Invalid encryption key or email"
// @Router /login [post]
// @Tags auth
// @x-order 3
func Login(c *gin.Context) {
	// Authenticate the user
	var req structs.LoginRequest
	user, err := Authenticate(c, &req)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key or email", false) {
		return
	}

	// Decrypt the Telegram bot API key and user ID
	decryptedBotAPIKey, err := encryption.Decrypt(env.MK, user.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	decryptedUserID, err := encryption.Decrypt(env.MK, user.TelegramUserID)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	// Load timezone
	location, err := time.LoadLocation(user.Timezone)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid timezone", false) {
		return
	}

	// Combine the TIME with a date to handle timezone conversion
	now := time.Now()
	rTAsTime, err := time.Parse("15:04:05", user.ReminderTime)
	reminderTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		rTAsTime.Hour(), rTAsTime.Minute(), rTAsTime.Second(), 0,
		time.UTC,
	)

	// Convert the combined time to the user's timezone
	reminderTimeLocal := reminderTime.In(location).Format("15:04")

	// Find the birthdays by user id
	birthdays, err := models.Birthdays(models.BirthdayWhere.UserID.EQ(user.ID.Int64), qm.Select("id", "name", "date")).All(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to fetch birthdays", false) {
		return
	}

	// Create a new slice to hold the filtered birthday data
	var filteredBirthdays []structs.BirthdayFull

	// Iterate over the birthdays and append the filtered data to the new slice
	for _, birthday := range birthdays {
		filteredBirthdays = append(filteredBirthdays, structs.BirthdayFull{
			ID:   birthday.ID.Int64,
			Name: birthday.Name,
			Date: birthday.Date.Format("2006-01-02"),
		})
	}

	// Return the user's details along with the filtered birthdays
	c.JSON(http.StatusOK, structs.LoginSuccess{
		TelegramBotAPIKey: decryptedBotAPIKey,
		TelegramUserID:    decryptedUserID,
		ReminderTime:      reminderTimeLocal,
		Timezone:          user.Timezone,
		Birthdays:         filteredBirthdays,
	})
}

// @Summary Modify a user's details
// @Description This endpoint modifies a user's details such as Telegram bot API key, reminder time, and more.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.ModifyUserRequest  true  "Modify user"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Invalid encryption key"
// @Failure 500 {object} structs.Error "Failed to update user"
// @Router /modify-user [put]
// @Tags auth
// @x-order 4
func ModifyUser(c *gin.Context) {
	// Authenticate the user
	var req structs.ModifyUserRequest
	user, err := Authenticate(c, &req)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key or email", false) {
		return
	}

	// Load the new timezone
	location, err := time.LoadLocation(req.NewTimezone)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid timezone", false) {
		return
	}

	// Parse the new reminder time
	now := time.Now()
	reminderTime, err := time.ParseInLocation("15:04", req.NewReminderTime, location)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid reminder time format", false) {
		return
	}

	// Encrypt the new Telegram bot API key and user ID
	telegramBotAPIKeyHash := encryption.HashStringWithSHA256(req.NewTelegramBotAPIKey)
	encryptedBotAPIKey, err := encryption.Encrypt(env.MK, req.NewTelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram bot API key", false) {
		return
	}

	telegramUserIDHash := encryption.HashStringWithSHA256(req.NewTelegramUserID)
	encryptedUserID, err := encryption.Encrypt(env.MK, req.NewTelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to encrypt Telegram user ID", false) {
		return
	}

	// Validate and hash the user's email (to be updated)
	if req.NewEmail != "" {
		if !helper.IsValidEmail(req.NewEmail) {
			c.JSON(http.StatusBadRequest, structs.Error{Error: "Invalid email format"})
			return
		}
		emailHash := encryption.HashStringWithSHA256(req.NewEmail)
		user.EmailHash = emailHash
	}

	// Ensure reminderTime is in UTC
	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).In(time.UTC)

	// Update the user's details
	user.ReminderTime = reminderTime.String()
	user.Timezone = req.NewTimezone
	user.TelegramBotAPIKey = hex.EncodeToString(encryptedBotAPIKey)
	user.TelegramBotAPIKeyHash = telegramBotAPIKeyHash
	user.TelegramUserID = hex.EncodeToString(encryptedUserID)
	user.TelegramUserIDHash = telegramUserIDHash

	// Start a new transaction
	tx, err := env.DB.Begin()
	if err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "Failed to begin transaction", false)
		return
	}

	// Perform the update within the transaction
	_, err = user.Update(c, tx, boil.Infer())
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		helper.HE(c, err, http.StatusInternalServerError, "Failed to update user", false)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "Failed to commit transaction", false)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// @Summary Delete a user
// @Description This endpoint deletes a user based on their email and encryption key.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.LoginRequest  true  "Delete user"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Invalid encryption key or email"
// @Failure 500 {object} structs.Error "Failed to delete user"
// @Router /delete-user [delete]
// @Tags auth
// @x-order 5
func DeleteUser(c *gin.Context) {
	// Authenticate the user
	var req structs.LoginRequest
	user, err := Authenticate(c, &req)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key or email", false) {
		return
	}

	// Start a new transaction
	_, err = user.Delete(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to delete user", false) {
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, structs.Success{Success: true})
}
