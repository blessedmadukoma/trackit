// File dealing with expense, budget, amount and income tables
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
)

// var jwtKey = []byte("my_secret_key")

// Dashboard
func Dashboard(r *http.Request) (models.User, models.ErrorResponse) {
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

	claimedUser, err := Dashboard(r)
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
	// if what is being returned is empty, i.e. result.BudgetName == "" && result.amount== "" && result.Description == "", return to the frontend: 0

	claimedUser, err := Dashboard(r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	budget := &models.Budget{}
	result := h.DB.Table("bugdets").First(&budget).Where("UserID", claimedUser.ID)
	if result.Error != nil {
		err := models.ErrorResponse{
			Message: `error getting a record`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(budget)
}

// Create Budget
func (h Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	// create budget
}

// Expense Screen

// Get Expense
func (h Handler) GetExpense(w http.ResponseWriter, r *http.Request) {
	// get expense from Expense table

	claimedUser, err := Dashboard(r)
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
	json.NewEncoder(w).Encode(expense)
}

// Add expense
func (h Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	// Enter amount, description, select date (purchased), and category gotten from JSON
}

// Balance
// Get Balance for Dashboard --
func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// get balance from account table

	claimedUser, err := Dashboard(r)
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
