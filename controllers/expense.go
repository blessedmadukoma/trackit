package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/blessedmadukoma/trackit-chima/models"
	"github.com/go-playground/validator/v10"
)

// Expense Screen

// Get all expenses
func (h Handler) GetAllExpenses(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	result, errur := h.DB.Raw(`SELECT * FROM expenses WHERE user_id=?`, claimedUser.ID).Rows()

	if errur != nil {
		err := models.ErrorResponse{
			Message: `error getting all expenses`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	if result == nil {
		w.WriteHeader(http.StatusNotFound)

		errResponse := &models.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "No expense found!",
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	expenses, expense := []models.Expense{}, models.Expense{}

	for result.Next() {
		err := result.Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt, &expense.DeletedAt, &expense.Amount, &expense.Description, &expense.Date_purchased, &expense.Category, &expense.UserID)
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
		expense.User.ID = claimedUser.ID
		expense.User.Firstname = claimedUser.Firstname
		expense.User.Lastname = claimedUser.Lastname
		expense.User.Email = claimedUser.Email
		expense.User.Mobile = claimedUser.Mobile
		expenses = append(expenses, expense)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
	defer result.Close()
}

// Get Expense
func (h Handler) GetExpense(w http.ResponseWriter, r *http.Request) {
	// get expense from Expense table

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	expense := &models.Expense{}
	result := h.DB.Table("expenses").First(&expense).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		errResponse := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	expense.User.ID = claimedUser.ID
	expense.User.Firstname = claimedUser.Firstname
	expense.User.Lastname = claimedUser.Lastname
	expense.User.Email = claimedUser.Email
	expense.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(expense)
}

// Add expense
func (h Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	// Enter amount, description, select date (purchased), and category gotten from JSON

	var Validator = validator.New()

	expense := &models.Expense{}
	json.NewDecoder(r.Body).Decode(&expense)

	validationError := Validator.Struct(expense)
	if validationError != nil {
		log.Fatal("Error validating struct:", validationError)
	}
	claimedUser, errorResponse := Dashboard(w, r)
	if errorResponse.Message != "" {
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if expense.Amount < 1 || expense.Description == "" {
		errorResponse := models.ErrorResponse{
			Message: `Invalid values`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// assign some values
	expense.User.ID = claimedUser.ID
	expense.User.Firstname = claimedUser.Firstname
	expense.User.Lastname = claimedUser.Lastname
	expense.User.Email = claimedUser.Email
	expense.User.Mobile = claimedUser.Mobile

	date := time.Now().UTC()
	time := fmt.Sprint(date.Hour()+1) + ":" + fmt.Sprint(date.Minute())

	transaction := &models.Transactions{}
	transaction.Amount = expense.Amount
	transaction.Category = expense.Category
	transaction.Date = expense.Date_purchased
	transaction.Time = time
	transaction.UserID = claimedUser.ID
	transaction.User.ID = claimedUser.ID
	transaction.User.Firstname = claimedUser.Firstname
	transaction.User.Lastname = claimedUser.Lastname
	transaction.User.Email = claimedUser.Email
	transaction.User.Mobile = claimedUser.Mobile

	result := h.DB.Create(&expense).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		fmt.Println("result error:", result.Error)
		errorResponse := models.ErrorResponse{
			Message: `error creating expense`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	errorr := h.DB.Create(&transaction).Where("user_id", claimedUser.ID)
	if errorr.Error != nil {
		err := models.ErrorResponse{
			Message: `error saving to transactions`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(expense)
}
