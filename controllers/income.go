package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blessedmadukoma/trackit-chima/models"
	"github.com/go-playground/validator/v10"
)

// Income screen
// Get all income
func (h Handler) GetAllIncome(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	result, errur := h.DB.Raw(`SELECT * FROM incomes WHERE user_id=$1`, claimedUser.ID).Rows()

	if errur != nil {
		errorResponse := models.ErrorResponse{
			Message: `error getting all incomes`,
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
			Message: "No income found!",
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	incomes, income := []models.Income{}, models.Income{}

	for result.Next() {
		err := result.Scan(&income.ID, &income.CreatedAt, &income.UpdatedAt, &income.DeletedAt, &income.Amount, &income.Date, &income.UserID)
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
		income.User.ID = claimedUser.ID
		income.User.Firstname = claimedUser.Firstname
		income.User.Lastname = claimedUser.Lastname
		income.User.Email = claimedUser.Email
		income.User.Mobile = claimedUser.Mobile
		incomes = append(incomes, income)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incomes)
	defer result.Close()
}

// Get Income --> update this route so that it is dynamic /income/{id}
func (h Handler) GetIncome(w http.ResponseWriter, r *http.Request) {

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	income := &models.Income{}
	result := h.DB.Table("incomes").First(&income).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	income.User.ID = claimedUser.ID
	income.User.Firstname = claimedUser.Firstname
	income.User.Lastname = claimedUser.Lastname
	income.User.Email = claimedUser.Email
	income.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(income)
}

// Add income screen
func (h Handler) AddIncome(w http.ResponseWriter, r *http.Request) {
	// add amount, date gotten from json
	var Validator = validator.New()

	income := &models.Income{}
	json.NewDecoder(r.Body).Decode(&income)

	validationError := Validator.Struct(income)
	if validationError != nil {
		err := models.ErrorResponse{
			Message: `Values could not be validated`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}
	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	// set income amount to be nothing less than 50
	if income.Amount < 50 {
		err := models.ErrorResponse{
			Message: `Income cannot be that low!`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	// date := income.Date

	// assign some values
	income.User.ID = claimedUser.ID
	income.User.Firstname = claimedUser.Firstname
	income.User.Lastname = claimedUser.Lastname
	income.User.Email = claimedUser.Email
	income.User.Mobile = claimedUser.Mobile

	date := time.Now().UTC()
	time := fmt.Sprint(date.Hour()+1) + ":" + fmt.Sprint(date.Minute())

	transaction := &models.Transactions{}
	transaction.Amount = income.Amount
	transaction.Category = "income"
	transaction.Date = income.Date
	transaction.Time = time
	transaction.UserID = claimedUser.ID
	transaction.User.ID = claimedUser.ID
	transaction.User.Firstname = claimedUser.Firstname
	transaction.User.Lastname = claimedUser.Lastname
	transaction.User.Email = claimedUser.Email
	transaction.User.Mobile = claimedUser.Mobile

	result := h.DB.Create(&income).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		// fmt.Println("result error:", result.Error)
		err := models.ErrorResponse{
			Message: `error creating income`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
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

	json.NewEncoder(w).Encode(income)
}
