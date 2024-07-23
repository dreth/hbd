package auth

import (
	"errors"
	"hbd/encryption"
	"hbd/env"
	"hbd/models"
	"hbd/structs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetUserByEmail retrieves a user from the database using their email address.
// The email is extracted from the Gin context, hashed using SHA-256, and the corresponding user is fetched.
//
// Parameters:
//   - c: The Gin context, which provides request and response handling.
//
// Returns:
//   - *models.User: The user object if found.
//   - error: An error object if any error occurs during the process, such as an invalid email or if the user is not found.
//
// The function performs the following steps:
// 1. Extracts the email from the Gin context.
// 2. Hashes the email using SHA-256.
// 3. Queries the database for a user with the given email hash.
// 4. Returns the user object or an error if the user is not found or an error occurs.
func GetUserByEmail(c *gin.Context) (*models.User, error) {
	// Get the email from the context
	email := c.GetString("Email")

	// Hash the email
	emailHash := encryption.HashStringWithSHA256(email)

	// Fetch the user with the given email hash from the database
	user, err := models.Users(
		qm.Where("email_hash = ?", emailHash),
	).One(c.Request.Context(), boil.GetContextDB())
	if err != nil {
		return nil, errors.New("invalid email")
	}

	return user, nil
}

// GetUserData fetches and returns user data including decrypted Telegram bot API key and user ID,
// reminder time converted to the user's timezone, and a list of birthdays associated with the user.
// It extracts the email from the Gin context, hashes it, and retrieves the corresponding user from the database.
//
// The function handles the following steps:
// 1. Extracts the email from the Gin context.
// 2. Hashes the email using SHA-256.
// 3. Queries the database for a user with the given email hash.
// 4. Decrypts the Telegram bot API key and user ID stored in the database.
// 5. Loads the user's timezone and parses the reminder time.
// 6. Converts the reminder time to the user's local timezone.
// 7. Retrieves and filters the list of birthdays associated with the user.
//
// Errors are returned if any step fails, including invalid email, decryption errors, invalid timezone,
// invalid reminder time format, and failure to fetch birthdays.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the request and user email.
//
// Returns:
// - (*structs.UserData, error): A pointer to the UserData struct containing user details and birthdays, or an error.
func GetUserData(c *gin.Context) (*structs.UserData, error) {
	// Get the user by its email
	user, err := GetUserByEmail(c)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	// Decrypt the Telegram bot API key and user ID
	decryptedBotAPIKey, err := encryption.Decrypt(env.MK, user.TelegramBotAPIKey)
	if err != nil {
		return nil, errors.New("error decrypting Telegram bot API key")
	}

	decryptedUserID, err := encryption.Decrypt(env.MK, user.TelegramUserID)
	if err != nil {
		return nil, errors.New("error decrypting Telegram user ID")
	}

	// Load timezone
	location, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return nil, errors.New("invalid timezone")
	}

	// Combine the TIME with a date to handle timezone conversion
	now := time.Now()
	rTAsTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", user.ReminderTime)
	if err != nil {
		return nil, errors.New("invalid reminder time format")
	}
	reminderTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		rTAsTime.Hour(), rTAsTime.Minute(), rTAsTime.Second(), 0,
		time.UTC,
	)

	// Convert the combined time to the user's timezone
	reminderTimeLocal := reminderTime.In(location).Format("15:04")

	// Find the birthdays by user id
	birthdays, err := models.Birthdays(models.BirthdayWhere.UserID.EQ(user.ID.Int64), qm.Select("id", "name", "date")).All(c, env.DB)
	if err != nil {
		return nil, errors.New("failed to fetch birthdays")
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

	// User data
	userData := structs.UserData{
		ID:                user.ID.Int64,
		TelegramBotAPIKey: decryptedBotAPIKey,
		TelegramUserID:    decryptedUserID,
		ReminderTime:      reminderTimeLocal,
		Timezone:          user.Timezone,
		Birthdays:         filteredBirthdays,
	}

	return &userData, nil
}
