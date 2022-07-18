package routes

import (
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/controllers"
	db "github.com/blessedmadukoma/trackit-chima/models"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {
	router := mux.NewRouter()

	DB := db.Init()
	h := controllers.New(DB)

	router.Use(CommonMiddleware)

	router.HandleFunc("/", h.Index).Methods("GET", "POST", "OPTIONS")

	// auth routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(CommonMiddleware)
	authRouter.HandleFunc("/signup", h.SignUp).Methods("GET", "POST", "OPTIONS")
	authRouter.HandleFunc("/signin", h.SignIn).Methods("POST", "GET", "OPTIONS")
	authRouter.HandleFunc("/logout", h.LogOut).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/reset-password", h.ResetPassword).Methods("PUT", "OPTIONS")

	// dashboard routes
	// balance
	router.HandleFunc("/balance", h.GetBalance).Methods("GET", "OPTIONS")
	// income
	router.HandleFunc("/income", h.GetIncome).Methods("GET", "OPTIONS")
	router.HandleFunc("/incomes", h.GetAllIncome).Methods("GET", "OPTIONS")
	router.HandleFunc("/income", h.AddIncome).Methods("POST", "OPTIONS")
	// expense
	router.HandleFunc("/expense", h.GetExpense).Methods("GET", "OPTIONS")
	router.HandleFunc("/expenses", h.GetAllExpenses).Methods("GET", "OPTIONS")
	router.HandleFunc("/expense", h.AddExpense).Methods("POST", "OPTIONS")
	// budget
	router.HandleFunc("/budget", h.GetBudget).Methods("GET", "OPTIONS")
	router.HandleFunc("/budget", h.UpdateBudget).Methods("PUT", "OPTIONS")
	// transactions
	router.HandleFunc("/transactions", h.GetTransactions).Methods("GET", "OPTIONS")

	// Savings
	router.HandleFunc("/savings", h.GetSavings).Methods("GET", "OPTIONS")

	// user routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(CommonMiddleware)
	// userRouter.Use(controllers.JwtVerify)
	// userRouter.HandleFunc("/dashboard", h.Dashboard)
	userRouter.HandleFunc("/users", h.FetchUsers).Methods("GET", "OPTIONS")
	userRouter.HandleFunc("/{id}", h.GetUser).Methods("GET", "OPTIONS")
	userRouter.HandleFunc("/{id}", h.UpdateUser).Methods("PUT", "OPTIONS")
	userRouter.HandleFunc("/{id}", h.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}

// --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
