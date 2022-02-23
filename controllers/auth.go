package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//SignUp function -- create a new user
func (h handler) SignUp(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
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
	json.NewEncoder(w).Encode(user)
}

// SignIn function
func (h handler) SignIn(w http.ResponseWriter, r *http.Request) {
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

	user := models.User{}

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
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &models.Claims{
			User: user,
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
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		// http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
		json.NewEncoder(w).Encode(user)
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
