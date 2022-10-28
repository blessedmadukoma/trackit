// Display transactions

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
)

func (h handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// display transactions

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	result, errur := h.DB.Raw(`SELECT * FROM transactions WHERE user_id=?`, claimedUser.ID).Rows()

	fmt.Println(errur)

	if errur != nil {
		errorResponse := models.ErrorResponse{
			Message: `error getting all transactions`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if result == nil {
		w.WriteHeader(http.StatusNotFound)

		errResponse := &models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "No transaction found!",
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	transactions, transaction := []models.Transactions{}, models.Transactions{}

	for result.Next() {
		err := result.Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt, &transaction.Category, &transaction.Amount, &transaction.Date, &transaction.Time, &transaction.UserID)
		if err != nil {
			fmt.Println(err)
			errResponse := &models.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "User not found!",
			}
			w.WriteHeader(errResponse.Status)
			json.NewEncoder(w).Encode(errResponse)
			return
		}
		if transaction.Category == "" && transaction.Date == "" {
			errResponse := &models.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "No transaction record!",
			}
			w.WriteHeader(errResponse.Status)
			json.NewEncoder(w).Encode(errResponse)
			return
		}
		transaction.User.ID = claimedUser.ID
		transaction.User.Firstname = claimedUser.Firstname
		transaction.User.Lastname = claimedUser.Lastname
		transaction.User.Email = claimedUser.Email
		transaction.User.Mobile = claimedUser.Mobile
		transaction.User.CreatedAt = claimedUser.CreatedAt
		transaction.User.UpdatedAt = claimedUser.UpdatedAt
		transactions = append(transactions, transaction)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
	defer result.Close()
}
