package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
)

// Balance
// Get Balance for Dashboard -- should be get all your balance -> you should have only one
func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// get balance from account table

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	account := &models.Account{}
	result := h.DB.Table("accounts").First(&account).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		errorResponse := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	account.User.ID = claimedUser.ID
	account.User.Firstname = claimedUser.Firstname
	account.User.Lastname = claimedUser.Lastname
	account.User.Email = claimedUser.Email
	account.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(account)
}

// Create Balance
// for every newly signed up user, they should have a balance of 0 by default, meaning the system autocreates a record for them
func (h Handler) CreateBalance(w http.ResponseWriter, r *http.Request) {
	// create balance - to be called at every successful sign up attempt
}

// Update Balance -- not sure this is useful because when they add income and expense, our balance will be affected and we will update the figure in that function either income or expense but the getBalance does the update since they are not manually updating it
func (h Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	// update balance - to be called at every transaction happening
}
