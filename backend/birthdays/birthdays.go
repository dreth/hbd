package birthdays

import (
	"hbd/auth"
	"hbd/encryption"
	"hbd/env"
	"hbd/helper"
	"hbd/models"
	"hbd/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// @Summary Check user reminders
// @Description This endpoint checks for user reminders through a POST request.
// @Accept  json
// @Produce  json
// @Param   user  body     structs.LoginRequest  true  "Check reminders"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid request"
// @Failure 500 {object} structs.Error "Error querying users"
// @Router /check-birthdays [post]
// @Tags reminders
// @x-order 6
func CallReminderChecker(c *gin.Context) {
	// Declare a variable to hold the login request data
	var req structs.LoginRequest

	// Authenticate the user and get user information
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
		return
	}

	// Decrypt the user's Telegram bot API key
	decryptedBotAPIKey, err := encryption.Decrypt(env.MK, user.TelegramBotAPIKey)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	// Decrypt the user's Telegram user ID
	decryptedUserID, err := encryption.Decrypt(env.MK, user.TelegramUserID)
	if helper.HE(c, err, http.StatusUnauthorized, "Invalid encryption key", false) {
		return
	}

	// Send the birthday reminder
	sendBirthdayReminder(user.ID, decryptedBotAPIKey, decryptedUserID)

	// Respond with a success message
	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// @Summary Add a new birthday
// @Description This endpoint adds a new birthday for the authenticated user.
// @Accept  json
// @Produce  json
// @Param   birthday  body     structs.BirthdayNameDateAdd  true  "Add birthday"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid encryption key or email"
// @Failure 400 {object} structs.Error "Invalid date format"
// @Failure 500 {object} structs.Error "Failed to insert birthday"
// @Router /birthday [post]
// @Tags birthdays
// @x-order 7
func AddBirthday(c *gin.Context) {
	// Declare a variable to hold the request data
	var req structs.BirthdayNameDateAdd

	// Authenticate the user and get user information
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid encryption key or email", true) {
		return
	}

	// Parse the date from the request
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helper.HE(c, err, http.StatusBadRequest, "Invalid date format", true)
		return
	}

	// Create a Birthday model with the parsed data
	b := models.Birthday{
		UserID: user.ID,
		Name:   req.Name,
		Date:   date,
	}

	// Insert the birthday into the database
	err = b.Insert(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to insert birthday", false) {
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// @Summary Delete a birthday
// @Description This endpoint deletes a birthday for the authenticated user.
// @Accept  json
// @Produce  json
// @Param   birthday  body     structs.BirthdayNameDateModify  true  "Delete birthday"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid encryption key or email"
// @Failure 400 {object} structs.Error "Invalid date format"
// @Failure 500 {object} structs.Error "Failed to delete birthday"
// @Router /birthday [delete]
// @Tags birthdays
// @x-order 8
func DeleteBirthday(c *gin.Context) {
	// Declare a variable to hold the request data
	var req structs.BirthdayNameDateModify

	// Authenticate the user and get user information
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid encryption key or email", true) {
		return
	}

	// Parse the date from the request
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helper.HE(c, err, http.StatusBadRequest, "Invalid date format", true)
		return
	}

	// Perform the delete operation on the birthdays matching the criteria
	_, err = models.Birthdays(
		models.BirthdayWhere.UserID.EQ(user.ID),
		models.BirthdayWhere.ID.EQ(req.ID),
		models.BirthdayWhere.Name.EQ(req.Name),
		models.BirthdayWhere.Date.EQ(date),
	).DeleteAll(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to delete birthday", false) {
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// @Summary Modify a birthday
// @Description This endpoint modifies a birthday for the authenticated user.
// @Accept  json
// @Produce  json
// @Param   birthday  body     structs.BirthdayNameDateModify  true  "Modify birthday"
// @Success 200 {object} structs.Success
// @Failure 400 {object} structs.Error "Invalid encryption key or email"
// @Failure 400 {object} structs.Error "Invalid date format"
// @Failure 500 {object} structs.Error "Birthday doesn't exist"
// @Failure 500 {object} structs.Error "Failed to begin transaction"
// @Failure 500 {object} structs.Error "Failed to update birthday"
// @Failure 500 {object} structs.Error "Failed to commit transaction"
// @Router /birthday [put]
// @Tags birthdays
// @x-order 9
func ModifyBirthday(c *gin.Context) {
	var req structs.BirthdayNameDateModify
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid encryption key or email", true) {
		return
	}

	// Perform the update operation
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helper.HE(c, err, http.StatusBadRequest, "Invalid date format", true)
		return
	}

	// Get the birthday
	birthday, err := models.Birthdays(
		models.BirthdayWhere.UserID.EQ(user.ID),
		models.BirthdayWhere.ID.EQ(req.ID),
	).One(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Birthday doesn't exist", false) {
		return
	}

	// Update the birthday
	birthday.Name = req.Name
	birthday.Date = date

	// Start a new transaction
	tx, err := env.DB.Begin()
	if err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "Failed to begin transaction", false)
		return
	}

	// Perform the update within the transaction
	_, err = birthday.Update(c, tx, boil.Infer())
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		helper.HE(c, err, http.StatusInternalServerError, "Failed to update birthday", false)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		helper.HE(c, err, http.StatusInternalServerError, "Failed to commit transaction", false)
		return
	}

	c.JSON(http.StatusOK, structs.Success{Success: true})
}
