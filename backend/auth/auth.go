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
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// @Summary Generate a new password
// @Description This endpoint generates a new password for the user.
// @Produce  json
// @Success 200 {object} structs.Password
// @Failure 500 {object} structs.Error "Failed to generate password"
// @Router /generate-password [get]
// @Tags auth
// @x-order 1
func GetPassword(c *gin.Context) {
	// Create a byte slice to hold the generated password
	key := make([]byte, 16)

	// Generate a random 16-byte key
	_, err := rand.Read(key)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to generate password", false) {
		return
	}

	// Respond with the generated password in hexadecimal format
	c.JSON(http.StatusOK, structs.Password{Password: hex.EncodeToString(key)})
}

// @Summary Register a new user
// @Description This endpoint registers a new user with their email, Telegram bot API key, and other details.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.RegisterRequest  true  "Register user"
// @Success 200 {object} structs.LoginSuccess
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

	// Check the length of the new email < 150 chars
	lengthErrors := helper.CheckArrayStringLength(
		[]string{"Email", "Password", "ReminderTime", "Timezone", "TelegramBotAPIKey", "TelegramUserID"},
		[]string{req.Email, req.Password, req.ReminderTime, req.Timezone, req.TelegramBotAPIKey, req.TelegramUserID},
		[]int{150, 64, 50, 50, 60, 20},
		[]int{5, 1, 1, 1, 1, 1},
		[]int{0, 0, 0, 0, 0, 0},
	)
	// Loop over errors and concatenate the strings to return it all at once
	// first check if all errors are nil, if so, nothing to do, otherwise, loop over the errors and concatenate them
	if helper.CheckErrors(lengthErrors) != nil {
		errorStr := helper.ConcatenateErrors(lengthErrors)
		c.JSON(http.StatusBadRequest, structs.Error{Error: errorStr})
		return
	}

	// Hash the email to check for uniqueness
	emailHash := encryption.HashStringWithSHA256(req.Email)

	// Check if the email hash already exists in the database
	exists, err := models.Users(models.UserWhere.EmailHash.EQ(emailHash)).Exists(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash the password to check for uniqueness
	PasswordHash := encryption.HashStringWithSHA256(req.Password)
	if exists, err = models.Users(models.UserWhere.PasswordHash.EQ(PasswordHash)).Exists(c, env.DB); helper.HE(c, err, http.StatusInternalServerError, "failed to check existing user", false) {
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "password already registered"})
		return
	}

	encryptedBotAPIKey, err := encryption.Encrypt(env.MK, req.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to encrypt Telegram bot API key", false) {
		return
	}

	encryptedUserID, err := encryption.Encrypt(env.MK, req.TelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to encrypt Telegram user ID", false) {
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

	// Create a new user object
	user := models.User{
		EmailHash:             emailHash,
		PasswordHash:          encryption.HashStringWithSHA256(req.Password),
		ReminderTime:          reminderTime.Format("15:04"),
		Timezone:              req.Timezone,
		TelegramBotAPIKey:     hex.EncodeToString(encryptedBotAPIKey),
		TelegramBotAPIKeyHash: encryption.HashStringWithSHA256(req.TelegramUserID),
		TelegramUserID:        hex.EncodeToString(encryptedUserID),
		TelegramUserIDHash:    encryption.HashStringWithSHA256(req.TelegramUserID),
	}

	err = user.Insert(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "failed to create user", false) {
		return
	}

	// As the user was successfully created, send a telegram message through the bot and ID to confirm the registration
	telegram.SendTelegramMessage(req.TelegramBotAPIKey, req.TelegramUserID, fmt.Sprintf("🎂 Your user has been successfully registered, through this bot and user ID you'll receive your birthday reminders (if there's any) at %s (Timezone: %s).\n\nIf you encounter any issues using the app or want to give any feedback to us. Please open an issue here: https://github.com/dreth/hbd/issues, thanks and we hope you find the application useful!", req.ReminderTime, req.Timezone))

	// Return the token and the user's details
	token, err := GenerateJWT(req.Email)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to generate token", false) {
		return
	} else {
		c.JSON(http.StatusOK, structs.LoginSuccess{
			Token:             token,
			TelegramBotAPIKey: req.TelegramBotAPIKey,
			TelegramUserID:    req.TelegramUserID,
			ReminderTime:      req.ReminderTime,
			Timezone:          req.Timezone,
			Birthdays:         []structs.BirthdayFull{},
		})
	}
}

// @Summary Login a user
// @Description This endpoint logs in a user by validating their email and password. Upon successful authentication, it generates a JWT token and returns the user's details along with the filtered list of birthdays.
//
// The login process includes the following steps:
// 1. Parses the login request payload to extract email and password.
// 2. Authenticates the user using the provided credentials.
// 3. Fetches the user data including decrypted Telegram bot API key, user ID, reminder time in local timezone, and birthdays.
// 4. Generates a JWT token for the authenticated user.
// 5. Returns the user's details, JWT token, and birthdays.
//
// Errors:
// - Returns 400 if the request payload is invalid.
// - Returns 401 if the authentication fails due to invalid email or password.
// - Returns 500 if there is an internal server error while fetching user data or generating the JWT token.
//
// @Accept  json
// @Produce  json
// @Param   user  body     structs.LoginRequest  true  "Login user"
// @Success 200 {object} structs.LoginSuccess
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Invalid email or password"
// @Failure 500 {object} structs.Error "Internal server error"
// @Router /login [post]
// @Tags auth
// @x-order 3
func Login(c *gin.Context) {
	// Parse the login request payload
	var req structs.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.Error{Error: "Invalid request"})
		return
	}

	// Check the length of the email and password
	lengthErrors := helper.CheckArrayStringLength(
		[]string{"Email", "Password"},
		[]string{req.Email, req.Password},
		[]int{150, 64},
		[]int{5, 1},
		[]int{0, 0},
	)
	if helper.CheckErrors(lengthErrors) != nil {
		errorStr := helper.ConcatenateErrors(lengthErrors)
		c.JSON(http.StatusBadRequest, structs.Error{Error: errorStr})
		return
	}

	// Hash the email and password
	emailHash := encryption.HashStringWithSHA256(req.Email)
	passwordHash := encryption.HashStringWithSHA256(req.Password)

	// Fetch the user with the given email hash and password hash from the database
	_, err := models.Users(
		qm.Where("email_hash = ?", emailHash),
		qm.Where("password_hash = ?", passwordHash),
	).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		c.JSON(http.StatusUnauthorized, structs.Error{Error: "invalid email or password"})
		return
	}

	// Set the user email in the context
	c.Set("Email", req.Email)

	userData, err := GetUserData(c)
	if helper.HE(c, err, http.StatusInternalServerError, "invalid email or password", true) {
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(req.Email)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to generate token", false) {
		return
	}

	// Return the user's details along with the filtered birthdays
	c.JSON(http.StatusOK, structs.LoginSuccess{
		Token:             token,
		TelegramBotAPIKey: userData.TelegramBotAPIKey,
		TelegramUserID:    userData.TelegramUserID,
		ReminderTime:      userData.ReminderTime,
		Timezone:          userData.Timezone,
		Birthdays:         userData.Birthdays,
	})
}

// @Summary Get user data
// @Description This endpoint returns the authenticated user's data including Telegram bot API key, user ID, reminder time, and birthdays. The request must include a valid JWT token.
// @Produce  json
// @Success 200 {object} structs.UserData
// @Failure 500 {object} structs.Error "Internal server error"
// @Security Bearer
// @Router /me [get]
// @Tags user
func Me(c *gin.Context) {
	userData, err := GetUserData(c)
	if helper.HE(c, err, http.StatusInternalServerError, "invalid email or password", true) {
		return
	}

	c.JSON(http.StatusOK, userData)
}

// @Summary Modify a user's details
// @Description This endpoint modifies a user's details such as Telegram bot API key, reminder time, and more. The request must include a valid JWT token.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.ModifyUserRequest  true  "Modify user"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Unauthorized"
// @Failure 500 {object} structs.Error "Failed to update user"
// @Security Bearer
// @Router /modify-user [put]
// @Tags auth
// @x-order 4
func ModifyUser(c *gin.Context) {
	// Retrieve the user from the database
	user, err := GetUserByEmail(c)
	if helper.HE(c, err, http.StatusUnauthorized, "invalid email", false) {
		return
	}

	// Parse the request body
	var req structs.ModifyUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.Error{Error: "Invalid request"})
		return
	}

	// Check the length of the new email and other fields
	lengthErrors := helper.CheckArrayStringLength(
		[]string{"NewEmail", "NewPassword", "NewReminderTime", "NewTimezone", "NewTelegramBotAPIKey", "NewTelegramUserID"},
		[]string{req.NewEmail, req.NewPassword, req.NewReminderTime, req.NewTimezone, req.NewTelegramBotAPIKey, req.NewTelegramUserID},
		[]int{150, 64, 50, 50, 60, 20},
		[]int{5, 5, 1, 1, 1, 1},
		[]int{0, 0, 0, 0, 0, 0},
	)
	if helper.CheckErrors(lengthErrors) != nil {
		errorStr := helper.ConcatenateErrors(lengthErrors)
		c.JSON(http.StatusBadRequest, structs.Error{Error: errorStr})
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
	if helper.HE(c, err, http.StatusInternalServerError, "failed to encrypt Telegram bot API key", false) {
		return
	}

	telegramUserIDHash := encryption.HashStringWithSHA256(req.NewTelegramUserID)
	encryptedUserID, err := encryption.Encrypt(env.MK, req.NewTelegramUserID)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to encrypt Telegram user ID", false) {
		return
	}

	// Validate and hash the user's new email (to be updated)
	if req.NewEmail != "" {
		if !helper.IsValidEmail(req.NewEmail) {
			c.JSON(http.StatusBadRequest, structs.Error{Error: "invalid email format"})
			return
		}
		emailHash := encryption.HashStringWithSHA256(req.NewEmail)
		user.EmailHash = emailHash
	}

	// Validate and hash the user's new password (to be updated)
	if req.NewPassword != "" {
		PasswordHash := encryption.HashStringWithSHA256(req.NewPassword)
		user.PasswordHash = PasswordHash
	}

	// Ensure reminderTime is in UTC
	reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, location).In(time.UTC)

	// Update the user's details
	user.ReminderTime = reminderTime.Format("15:04")
	user.Timezone = req.NewTimezone
	user.TelegramBotAPIKey = hex.EncodeToString(encryptedBotAPIKey)
	user.TelegramBotAPIKeyHash = telegramBotAPIKeyHash
	user.TelegramUserID = hex.EncodeToString(encryptedUserID)
	user.TelegramUserIDHash = telegramUserIDHash

	// Start a new transaction
	tx, err := env.DB.Begin()
	if err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "failed to begin transaction", false)
		return
	}

	// Perform the update within the transaction
	_, err = user.Update(c, tx, boil.Infer())
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		helper.HE(c, err, http.StatusInternalServerError, "failed to update user", false)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "failed to commit transaction", false)
		return
	}

	// After committing the transaction, emit another JWT token with the new email
	token, err := GenerateJWT(req.NewEmail)
	if helper.HE(c, err, http.StatusInternalServerError, "failed to generate token", false) {
		return
	}

	// Update the email in the context
	c.Set("Email", req.NewEmail)

	// Get the user data post-changes
	userData, err := GetUserData(c)
	if helper.HE(c, err, http.StatusInternalServerError, "invalid email or password", true) {
		return
	}

	// Return the new token with the new user data
	c.JSON(http.StatusOK, structs.LoginSuccess{
		Token:             token,
		TelegramBotAPIKey: userData.TelegramBotAPIKey,
		TelegramUserID:    userData.TelegramUserID,
		ReminderTime:      userData.ReminderTime,
		Timezone:          userData.Timezone,
		Birthdays:         userData.Birthdays,
	})
}

// @Summary Delete a user
// @Description This endpoint deletes a user based on their email obtained from the JWT token. The request must include a valid JWT token.
// @Accept  json
// @Produce  json
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 401 {object} structs.Error "Unauthorized"
// @Failure 500 {object} structs.Error "Failed to delete user"
// @Security Bearer
// @Router /delete-user [delete]
// @Tags auth
// @x-order 5
func DeleteUser(c *gin.Context) {
	// Retrieve the user from the database
	user, err := GetUserByEmail(c)
	if helper.HE(c, err, http.StatusUnauthorized, "invalid email", false) {
		return
	}

	// Retrieve the user data and get the telegram bot api key and ID to send a confirmation message that the account has been deleted
	userData, err := GetUserData(c)
	if helper.HE(c, err, http.StatusInternalServerError, "invalid email or password", true) {
		return
	}

	// Start a new transaction
	tx, err := env.DB.Begin()
	if err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "failed to begin transaction", false)
		return
	}

	// Perform the delete within the transaction
	_, err = user.Delete(c, tx)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		helper.HE(c, err, http.StatusInternalServerError, "failed to delete user", false)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "failed to commit transaction", false)
		return
	}

	// Because the transaction is successful, let's retain
	tgbotapi.NewBotAPI(userData.TelegramBotAPIKey)
	telegram.SendTelegramMessage(userData.TelegramBotAPIKey, userData.TelegramUserID, "🎂 Your account and all your data has successfully been deleted forever. We're sorry to see you go ):\n\nThanks for checking out the app! If you have any feedback, feel free to open an issue: https://github.com/dreth/hbd/issues, we really appreciate it!")

	// Return a success response
	c.JSON(http.StatusOK, structs.Success{Success: true})
}
