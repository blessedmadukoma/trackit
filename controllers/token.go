// File dealing with expense, budget, amount and income tables
package controllers

import (
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
	// "gopkg.in/dealancer/validate.v2"
)

// var jwtKey = []byte("my_secret_key")

func GetJWT() []byte {
	jwtKey := []byte("qwerewndeinn3#456%$%^Y4&")
	return jwtKey
}

// Get user token if signed in
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
