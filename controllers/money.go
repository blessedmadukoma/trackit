// File dealing with expense, budget, amount and income tables
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	// "gopkg.in/dealancer/validate.v2"
)

// var jwtKey = []byte("my_secret_key")

// Dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) (models.User, models.ErrorResponse) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			err := models.ErrorResponse{
				Message: `No cookie set error`,
				Status:  http.StatusUnauthorized,
			}
			// w.WriteHeader(http.StatusUnauthorized)
			// json.NewEncoder(w).Encode(err)
			return models.User{}, err
		}
		// If the cookie is not set, return an unauthorized status
		err := models.ErrorResponse{
			Message: `Any other cookie error`,
			Status:  http.StatusUnauthorized,
		}
		// w.WriteHeader(http.StatusUnauthorized)
		// json.NewEncoder(w).Encode(err)
		return models.User{}, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &models.Claims{}

	jwtKey := GetJWT()

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			err := models.ErrorResponse{
				Message: `Invalid signatrue`,
				Status:  http.StatusUnauthorized,
			}
			// w.WriteHeader(http.StatusUnauthorized)
			// json.NewEncoder(w).Encode(err)
			return models.User{}, err
		}
		err := models.ErrorResponse{
			Message: `Error returning jwtKey`,
			Status:  http.StatusBadRequest,
		}
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(err)
		return models.User{}, err
	}
	if !tkn.Valid {
		err := models.ErrorResponse{
			Message: `token not valid`,
			Status:  http.StatusUnauthorized,
		}
		// w.WriteHeader(http.StatusUnauthorized)
		// json.NewEncoder(w).Encode(err)
		return models.User{}, err
	}

	// json.NewEncoder(w).Encode(claims.User)
	return claims.User, models.ErrorResponse{}
}

// Income screen
// Get Income
func (h Handler) GetIncome(w http.ResponseWriter, r *http.Request) {

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	income := &models.Income{}
	result := h.DB.Table("incomes").First(&income).Where("UserID", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	income.UserID = claimedUser.ID
	income.User.Firstname = claimedUser.Firstname
	income.User.Lastname = claimedUser.Lastname
	income.User.Email = claimedUser.Email
	income.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(income)
}

// Add income screen
func (h Handler) AddIncome(w http.ResponseWriter, r *http.Request) {
	// add amount, date gotten from json

	// initiailly set to Zero - remove when you have connected to database
	value := 0
	json.NewEncoder(w).Encode(value)
}

// func getIncome(id string) {
// 	h.DB.First(id)
// }

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

	budget.UserID = claimedUser.ID
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

	budget.UserID = claimedUser.ID
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

// Expense Screen

// Get Expense
func (h Handler) GetExpense(w http.ResponseWriter, r *http.Request) {
	// get expense from Expense table

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	// resultValue := h.DB.Raw("SELECT * FROM accounts LIMIT 1").Where("userID", claimedUser.ID)

	expense := &models.Expense{}
	result := h.DB.Table("expenses").First(&expense).Where("UserID", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	expense.UserID = claimedUser.ID
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
	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	if expense.Amount < 1 || expense.Description == "" {
		err := models.ErrorResponse{
			Message: `Invalid values`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	// assign some values
	expense.UserID = claimedUser.ID
	expense.User.Firstname = claimedUser.Firstname
	expense.User.Lastname = claimedUser.Lastname
	expense.User.Email = claimedUser.Email
	expense.User.Mobile = claimedUser.Mobile

	result := h.DB.Create(&expense).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		fmt.Println("result error:", result.Error)
		err := models.ErrorResponse{
			Message: `error creating expense`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(expense)
}

// Balance
// Get Balance for Dashboard --
func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// get balance from account table

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	account := &models.Account{}
	result := h.DB.Table("accounts").First(&account).Where("UserID", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	account.UserID = claimedUser.ID
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
