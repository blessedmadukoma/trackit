// Display transactions

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
)

func (h Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// display transactions
	// get the transactions: time or date, category and amount
	// if what is being returned is empty i.e. result.Amount == "" && result.Category == "", return to the frontend: No recent activities

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	result, errur := h.DB.Raw(`SELECT * FROM transactions WHERE user_id=?`, claimedUser.ID).Rows()

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

	// ID        uint `gorm:"primarykey"`
	//    CreatedAt time.Time
	//    UpdatedAt time.Time
	//    DeletedAt DeletedAt `gorm:"index"`

	// 			Category string  `json:"category" gorm:"not null"`
	// Amount   float64 `json:"amount" gorm:"not null"`
	// Date string `json:"date" gorm:"not null"`
	// Time string `json:"time" gorm:"not null"`
	// UserID   uint
	// User     User

	for result.Next() {
		err := result.Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt, &transaction.Category, &transaction.Amount, &transaction.UserID, &transaction.Date, &transaction.Time)
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
