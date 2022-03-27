package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/blessedmadukoma/trackit-chima/models"
)

func (h Handler) GetSavings(w http.ResponseWriter, r *http.Request) {
	// get budget

	claimedUser, err := Dashboard(w, r)
	if err.Message != "" {
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	savings := &models.Savings{}
	result := h.DB.Table("savings").First(&savings).Where("user_id", claimedUser.ID)
	if result.Error != nil {
		errResponse := models.ErrorResponse{
			Message: `no savings recorded. Try creating a budget!`,
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	savings.User.ID = claimedUser.ID
	savings.User.Firstname = claimedUser.Firstname
	savings.User.Lastname = claimedUser.Lastname
	savings.User.Email = claimedUser.Email
	savings.User.Mobile = claimedUser.Mobile

	json.NewEncoder(w).Encode(savings)
}
