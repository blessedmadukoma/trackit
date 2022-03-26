package routes

import (
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/controllers"
	db "github.com/blessedmadukoma/trackit-chima/models"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.Use(CommonMiddleware)

	DB := db.Init()
	h := controllers.New(DB)

	router.HandleFunc("/index", h.Index)

	// auth routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", h.SignUp).Methods("POST", "GET")
	authRouter.HandleFunc("/signin", h.SignIn).Methods("POST", "GET")
	authRouter.HandleFunc("/logout", h.LogOut).Methods("POST")
	authRouter.HandleFunc("/reset-password", h.ResetPassword).Methods("PUT")

// dashboard routes
// balance
	router.HandleFunc("/balance", h.GetBalance).Methods("GET")
// income
	router.HandleFunc("/income", h.GetIncome).Methods("GET")
	router.HandleFunc("/incomes", h.GetAllIncome).Methods("GET")
	router.HandleFunc("/income", h.AddIncome).Methods("POST")
	// expense
	router.HandleFunc("/expense", h.GetExpense).Methods("GET")
	router.HandleFunc("/expenses", h.GetAllExpenses).Methods("GET")
	router.HandleFunc("/expense", h.AddExpense).Methods("POST")
	// budget
	router.HandleFunc("/budget", h.GetBudget).Methods("GET")
	router.HandleFunc("/budget", h.UpdateBudget).Methods("PUT")
	// transactions
	router.HandleFunc("/transactions", h.GetTransactions).Methods("GET")

	// user routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(controllers.JwtVerify)
	// userRouter.HandleFunc("/dashboard", h.Dashboard)
	userRouter.HandleFunc("/users", h.FetchUsers).Methods("GET")
	userRouter.HandleFunc("/{id}", h.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id}", h.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", h.DeleteUser).Methods("DELETE")

	return router
}

// --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
