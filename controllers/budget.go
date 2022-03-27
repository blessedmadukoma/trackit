package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
)

// Budget Screen
// Get Budget
func (h Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	// get budget

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	budget := &models.Budget{}
	result := h.DB.Table("budgets").First(&budget).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	budget.User.ID = claimedUser.ID
	budget.User.Firstname = claimedUser.Firstname
	budget.User.Lastname = claimedUser.Lastname
	budget.User.Email = claimedUser.Email
	budget.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(budget)
}

// Update Budget - this is the only function working: since they have only one budget
func (h Handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}
	budgetInput := &models.Budget{}
	json.NewDecoder(r.Body).Decode(&budgetInput)

	budget := &models.Budget{}
	result := h.DB.Table("budgets").First(&budget).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		errResponse := models.ErrorResponse{
			Message: `error getting budget`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	if budgetInput.InitialAmount < float64(10000) {
		errorResponse := models.ErrorResponse{
			Message: `Budget is lower than 10000.`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// get balance and check if the budget amount is greater than the balance
	balance := &models.Account{}
	result = h.DB.Table("accounts").First(&balance).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		errResponse := models.ErrorResponse{
			Message: `error getting accounts for user`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	if budgetInput.InitialAmount > balance.Amount {
		errResponse := models.ErrorResponse{
			Message: `Your budget cannot be greater than balance. Multiply your income :)`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	balance.Amount = balance.Amount - budgetInput.InitialAmount
	balResult := h.DB.Save(&balance)
	if balResult.Error != nil {
		errResponse := &models.ErrorResponse{
			Message: `error updating account balance`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	savings := &models.Savings{}
	savingsDB := h.DB.Table("accounts").Where("user_id", claimedUser.ID).Find(&savings)
	if savingsDB.Error != nil {
		errorResponse := models.ErrorResponse{
			Message: `error getting savings!`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if savings.Amount == 0 {
		savings.Amount = (0.05 * budgetInput.InitialAmount)
		savings.UserID = claimedUser.ID
	
		createSavings := h.DB.Create(&savings).Where("user_id", claimedUser.ID)
		if createSavings.Error != nil {
			errorResponse := models.ErrorResponse{
				Message: `error creating savings`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	savings.Amount = (0.05 * budgetInput.InitialAmount)
	savings.UserID = claimedUser.ID

	createSavings := h.DB.Save(&savings)
	if createSavings.Error != nil {
		errorResponse := models.ErrorResponse{
			Message: `error updating savings`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	moneyNotSaved := (budgetInput.InitialAmount - (0.05 * budgetInput.InitialAmount))

	budget.Budget_name = budgetInput.Budget_name
	budget.InitialAmount = moneyNotSaved
	budget.CurrentAmount = moneyNotSaved // creates the budget with the same amount as the money that excludes the savings
	budget.Description = budgetInput.Description
	budget.StartDate = budgetInput.StartDate
	budget.EndDate = budgetInput.EndDate

	budget.User.ID = claimedUser.ID
	budget.User.Firstname = claimedUser.Firstname
	budget.User.Lastname = claimedUser.Lastname
	budget.User.Email = claimedUser.Email
	budget.User.Mobile = claimedUser.Mobile

	result = h.DB.Save(&budget)
	if result.Error != nil {
		err := &models.ErrorResponse{
			Message: `error updating budget`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(budget)
}
