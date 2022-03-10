// File dealing with expense, budget, amount and income tables
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Income screen
// Get Income
func (h Handler) GetIncome(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	value := id
	// value := h.DB.First(id)
	
	// value := getIncome(id)
	json.NewEncoder(w).Encode(value)
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
}



// Expense Screen

// Get Expense
func (h Handler) GetExpense(w http.ResponseWriter, r *http.Request) {
	// get expense from Expense table
	// if what is being returned is empty, i.e. result.Description == "" && result.Amount== "", return to the frontend: 0
}
// Add expense
func (h Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	// Enter amount, description, select date (purchased), and category gotten from JSON
}

// Balance
// Get Balance for Dashboard
func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// get balance from amount table
	// if what is being returned is empty, i.e. result.UserID == "" && result.amount== "" && result.Description == "", send a message to the frontend saying: No recent activities
}
// Create Balance
// for every newly signed up user, they should have a balance of 0 by default, meaning the system autocreates a record for them
func (h Handler) CreateBalance(w http.ResponseWriter, r *http.Request) {
	// create balance - to be called at every successful sign up attempt
}
// Update Balance
func (h Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	// update balance - to be called at every transaction happening
}