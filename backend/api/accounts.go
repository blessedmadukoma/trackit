package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/blessedmadukoma/trackit-chima/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID  int64   `json:"user_id" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

// createAccount creates a new account
func (srv *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("input not valid", err))
		return
	}

	arg := db.CreateAccountParams{
		UserID:  req.UserID,
		Balance: req.Balance,
	}

	account, err := srv.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("unable to create account", err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
	return
}

type getAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccountByID retrieves an existing account using the ID
func (srv *Server) getAccountByID(ctx *gin.Context) {
	var req getAccountByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("ID not valid", err))
		return
	}

	account, err := srv.store.GetAccountByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("No record of this account", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving account", err))
		return
	}

	fmt.Println(account)

	ctx.JSON(http.StatusOK, account)
	return
}

type getAccountByUserIDRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

// getAccountByUserID retrieves an existing account using the ID
func (srv *Server) getAccountByUserID(ctx *gin.Context) {
	var req getAccountByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("ID not valid", err))
		return
	}

	account, err := srv.store.GetAccountByUserID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("No record of this account by this user ID", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving account from this user ID", err))
		return
	}

	ctx.JSON(http.StatusOK, account)
	return
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts displays all accounts stored in the database
func (srv *Server) listAccounts(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid input", err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := srv.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error getting user accounts", err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
	return
}

type updateAccountRequestID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountRequestJSON struct {
	// UserID  int64   `json:"user_id" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

// updateAccount updates an account based on the specific ID
func (srv *Server) updateAccount(ctx *gin.Context) {
	var req1 updateAccountRequestID
	if err := ctx.ShouldBindUri(&req1); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid ID input", err))
		return
	}

	var req2 updateAccountRequestJSON

	if err := ctx.ShouldBindJSON(&req2); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid JSON param input", err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req1.ID,
		Balance: req2.Balance,
	}

	updatedUserAccount, err := srv.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error updating account", err))
		return
	}

	ctx.JSON(http.StatusOK, updatedUserAccount)
	return
}
