package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var user = &models.User{}

//SignUp function -- create a new user
func (h handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode("Signup Screen")
	} else if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&user)

		err := validate.Struct(user)
		if err != nil {
			err := models.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusUnauthorized,
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		// check if user exists
		err = h.checkExistingUser(user)
		if err != nil {
			err := models.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			err := models.ErrorResponse{
				Message: "Error generating hash for password",
				Status:  http.StatusBadRequest,
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		user.Password = string(pass)

		createdUser := h.DB.Create(&user)

		if createdUser.Error != nil {
			err := models.ErrorResponse{
				Message: `Error creating user`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		// Create an empty account
		account := &models.Account{}
		account.Amount = 0
		account.User = *user
		account.UserID = user.ID
		accountCreated := h.DB.Create(&account)
		if accountCreated.Error != nil {
			err := models.ErrorResponse{
				Message: `Error creating empty account`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		// account.User.Firstname = user.Firstname
		// account.User.Lastname = user.Lastname
		// account.User.Email = user.Email
		// account.User.Mobile = user.Mobile
		// account.User.Password = user.Password
		// h.DB.Save(&account)

		// create empty income
		income := &models.Income{}
		income.Amount = 0
		income.User = *user
		income.UserID = user.ID
		incomeCreated := h.DB.Create(&income)
		if incomeCreated.Error != nil {
			err := models.ErrorResponse{
				Message: `Error creating empty income`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		// income.User.Firstname = user.Firstname
		// income.User.Lastname = user.Lastname
		// income.User.Email = user.Email
		// income.User.Mobile = user.Mobile
		// income.User.Password = user.Password
		// h.DB.Save(&income)

		// create empty expense
		expense := &models.Expense{}
		expense.Amount = 0
		expense.User = *user
		expense.UserID = user.ID
		expenseCreated := h.DB.Create(&expense)
		if expenseCreated.Error != nil {
			err := models.ErrorResponse{
				Message: `Error creating empty expense`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		// expense.User.Firstname = user.Firstname
		// expense.User.Lastname = user.Lastname
		// expense.User.Email = user.Email
		// expense.User.Mobile = user.Mobile
		// expense.User.Password = user.Password
		// h.DB.Save(&expense)

		// create empty budget
		budget := &models.Budget{}
		budget.Amount = 0
		budget.User = *user
		budget.UserID = user.ID
		budgetCreated := h.DB.Create(&budget)
		if budgetCreated.Error != nil {
			err := models.ErrorResponse{
				Message: `Error creating empty budget`,
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		// budget.User.Firstname = user.Firstname
		// budget.User.Lastname = user.Lastname
		// budget.User.Email = user.Email
		// budget.User.Mobile = user.Mobile
		// budget.User.Password = user.Password
		// h.DB.Save(&budget)

		json.NewEncoder(w).Encode(user)
	}
}

// ----- SignIn function
func (h handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode("Login Screen")
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
		}

		var formatttedBody models.LoginUser
		err = json.Unmarshal(body, &formatttedBody)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = validate.Struct(formatttedBody)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		email, password := formatttedBody.Email, formatttedBody.Password
		emailBlank, passwordBlank := strings.Trim(email, " ") == "", strings.Trim(password, " ") == ""
		if emailBlank || passwordBlank {
			if email == "" && password == "" {
				err := models.ErrorResponse{
					Message: "Both email and password fields are empty!",
					Status:  http.StatusBadRequest,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			} else if email == "" {
				err := models.ErrorResponse{
					Message: "Email field is empty!",
					Status:  http.StatusBadRequest,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			} else {
				err := models.ErrorResponse{
					Message: "Password field is empty!",
					Status:  http.StatusBadRequest,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
		}

		// user := models.User{}

		// Check if user exists in the database
		result := h.DB.Where("email=?", email).Find(&user)
		if result.Error != nil {
			err := models.ErrorResponse{
				Message: "No user exists",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			err := models.ErrorResponse{
				Message: "Password does not match!",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		} else {

			// expiration time = 20 minutes
			expirationTime := time.Now().Add(20 * time.Minute)
			// Create the JWT claims, which includes the username and expiry time
			claims := &models.Claims{
				User: *user,
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}

			// Declare the token with the algorithm used for signing, and the claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Create the JWT string
			jwtKey := GetJWT()
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				log.Println("Error creating JWT return", err)
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// set the client cookie for "token" as the JWT generated, set an expiry time same as the token itself
			cookie := &http.Cookie{
				Name:    "token",
				Path:    "/",
				Value:   tokenString,
				Expires: expirationTime,
			}

			// cookie.Value = tokenString
			http.SetCookie(w, cookie)
			// http.SetCookie(w, &http.Cookie{
			// 	Name:    "token",
			// 	Value:   tokenString,
			// 	Expires: expirationTime,
			// })

			fmt.Println(cookie.Value)

			json.NewEncoder(w).Encode(user)
		}
	}
}

// logout function
func (h handler) LogOut(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		err := HandleError(err)
		json.NewEncoder(w).Encode(err)
		return
	}
	d := http.Cookie{
		Name:   c.Name,
		MaxAge: -1} // setting the maxAge < 0 deletes the cookie
	http.SetCookie(w, &d)
	http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)

}

func (h handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	newResetUser := &models.ResetUser{}
	json.NewDecoder(r.Body).Decode(&newResetUser)

	// Validate form input
	if strings.Trim(newResetUser.Email, " ") == "" || strings.Trim(newResetUser.New_password, "") == "" {
		errValue := errors.New("parameters cannot be empty")
		err := HandleError(errValue)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	if len(newResetUser.New_password) < 6 {
		errValue := errors.New("password length is short, should be more than 6")
		err := HandleError(errValue)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	// Check if user exists in the database
	result := h.DB.Table("users").Where("email=?", newResetUser.Email).Find(&newResetUser)
	if result.Error != nil {
		err := HandleError(result.Error)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	new_password, err := bcrypt.GenerateFromPassword([]byte(newResetUser.New_password), bcrypt.DefaultCost)
	if err != nil {
		err := HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	newResetUser.New_password = string(new_password)
	updatedUser := h.DB.Table("users").Where("email=?", newResetUser.Email).Update("password", newResetUser.New_password)

	if updatedUser.Error != nil {
		err := HandleError(updatedUser.Error)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(newResetUser)
}
