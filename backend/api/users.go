package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/blessedmadukoma/trackit-chima/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserAccountRequest struct {
	Firstname string `json:"firstname" binding:"required,min=3"`
	Lastname  string `json:"lastname" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,min=3"`
	Mobile    string `json:"mobile" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=3"`
}

// createUserAccount creates a new user account
func (srv *Server) createUserAccount(ctx *gin.Context) {
	var req createUserAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input info not valid", err))
		return
	}

	arg := db.CreateUserAccountParams{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Password:  req.Password,
	}

	userAccount, err := srv.store.CreateUserAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("unable to create user account", err))
		return
	}

	ctx.JSON(http.StatusCreated, userAccount)
	return
}

// deleteUserAccountRequest contains the params needed as input
type deleteUserAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteUserAccount deletes an existing user account using the ID
func (srv *Server) deleteUserAccount(ctx *gin.Context) {
	var req deleteUserAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("ID not valid", err))
		return
	}

	// retreive info
	_, err := srv.store.GetUserAccountByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse("Error retrieving user", err))
		return
	}

	err = srv.store.DeleteUserAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error deleting user", err))
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("User %d successfully deleted!", req.ID))
	return
}

type getUserAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getUserAccountByID retrieves an existing user account using the ID
func (srv *Server) getUserAccountByID(ctx *gin.Context) {
	var req getUserAccountByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("ID not valid", err))
		return
	}

	userAccount, err := srv.store.GetUserAccountByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("No record of this user", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving user account", err))
		return
	}

	fmt.Println(userAccount)

	ctx.JSON(http.StatusOK, userAccount)
	return
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listUsers displays all user accounts stored in the database
func (srv *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid input", err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userAccounts, err := srv.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error getting user accounts", err))
		return
	}

	ctx.JSON(http.StatusOK, userAccounts)
	return
}

type updateUserAccountRequestID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateUserAccountRequestJSON struct {
	Firstname string `json:"firstname" binding:"required,min=3"`
	Lastname  string `json:"lastname" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,min=3"`
	Mobile    string `json:"mobile" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=3"`
}

// updateUserAccount updates a user account based on the specific ID
func (srv *Server) updateUserAccount(ctx *gin.Context) {
	var req1 updateUserAccountRequestID
	if err := ctx.ShouldBindUri(&req1); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid input for ID", err))
		return
	}

	var req2 updateUserAccountRequestJSON

	if err := ctx.ShouldBindJSON(&req2); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid input for JSON parameters", err))
		return
	}

	arg := db.UpdateUserAccountParams{
		ID:        req1.ID,
		Firstname: req2.Firstname,
		Lastname:  req2.Lastname,
		Email:     req2.Email,
		Mobile:    req2.Mobile,
		Password:  req2.Password,
	}

	updatedUserAccount, err := srv.store.UpdateUserAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error updating user account", err))
		return
	}

	ctx.JSON(http.StatusOK, updatedUserAccount)
	return
}
