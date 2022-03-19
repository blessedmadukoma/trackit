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
// Get all income
func (h Handler) GetAllIncome(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	result, errur := h.DB.Raw(`SELECT * FROM incomes WHERE user_id=?`, claimedUser.ID).Rows()

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
		err := result.Scan(&income.ID, &income.CreatedAt, &income.UpdatedAt, &income.DeletedAt, &income.Amount, &income.Date,&income.UserID)
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

	// assign some values
	income.User.ID = claimedUser.ID
	income.User.Firstname = claimedUser.Firstname
	income.User.Lastname = claimedUser.Lastname
	income.User.Email = claimedUser.Email
	income.User.Mobile = claimedUser.Mobile

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
	json.NewEncoder(w).Encode(income)
}

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

// Expense Screen

// Get all expenses
// Get all income
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
	json.NewEncoder(w).Encode(expense)
}

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
