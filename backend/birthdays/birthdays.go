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
	var req structs.LoginRequest

	// Authenticate the user
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid request", true) {
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

	sendBirthdayReminder(user.ID, decryptedBotAPIKey, decryptedUserID)

	c.JSON(http.StatusOK, structs.Success{Success: true})
}

func AddBirthday(c *gin.Context) {
	var req structs.BirthdayNameDateAdd
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid encryption key or email", true) {
		return
	}

	// Perform the insert operation
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helper.HE(c, err, http.StatusBadRequest, "Invalid date format", true)
		return
	}

	b := models.Birthday{
		UserID: user.ID,
		Name:   req.Name,
		Date:   date,
	}

	err = b.Insert(c, env.DB, boil.Infer())
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to insert birthday", false) {
		return
	}

	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// Delete birthday func
func DeleteBirthday(c *gin.Context) {
	var req structs.BirthdayNameDateModify
	user, err := auth.Authenticate(c, &req)
	if helper.HE(c, err, http.StatusBadRequest, "Invalid encryption key or email", true) {
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		helper.HE(c, err, http.StatusBadRequest, "Invalid date format", true)
		return
	}

	// Perform the delete operation
	_, err = models.Birthdays(
		models.BirthdayWhere.UserID.EQ(user.ID),
		models.BirthdayWhere.ID.EQ(req.ID),
		models.BirthdayWhere.Name.EQ(req.Name),
		models.BirthdayWhere.Date.EQ(date),
	).DeleteAll(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Failed to delete birthday", false) {
		return
	}

	c.JSON(http.StatusOK, structs.Success{Success: true})
}

// Modify birthday func
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
		models.BirthdayWhere.Name.EQ(req.Name),
		models.BirthdayWhere.Date.EQ(date),
	).One(c, env.DB)
	if helper.HE(c, err, http.StatusInternalServerError, "Birthday doesn't exist", false) {
		return
	}

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
