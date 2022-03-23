// Display transactions

package controllers

import "net/http"

func GetTransactions(w http.ResponseWriter, r *http.Request) {
		// display transactions
		// get the transactions: time or date, category and amount
		// if what is being returned is empty i.e. result.Amount == "" && result.Category == "", return to the frontend: No recent activities
}
