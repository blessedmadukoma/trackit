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
		err := models.ErrorResponse{
			Message: `error getting budget`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	budget.Budget_name = budgetInput.Budget_name
	budget.Amount = budgetInput.Amount
	budget.Description = budgetInput.Budget_name
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
