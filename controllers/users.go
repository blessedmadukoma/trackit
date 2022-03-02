package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
	"github.com/go-playground/validator"

	"github.com/gorilla/mux"
)

var validate = validator.New()

// welcome function
func (h handler) Index(w http.ResponseWriter, r *http.Request) {

	value := "Index Page"
	json.NewEncoder(w).Encode(value)
}

func (h handler) FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	h.DB.Preload("auths").Find(&users)

	json.NewEncoder(w).Encode(users)
}

func (h handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	h.DB.First(&user, id)
	json.NewDecoder(r.Body).Decode(user)
	h.DB.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

func (h handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	h.DB.First(&user, id)
	h.DB.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}

func (h handler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	h.DB.First(&user, id)
	json.NewEncoder(w).Encode(&user)
}
