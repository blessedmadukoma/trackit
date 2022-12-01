package api

import (
	db "github.com/blessedmadukoma/trackit-chima/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server struct holds the configuration for the server db and router
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// proxies settings

	// Routes
	// user routes
	router.POST("/api/users", server.createUserAccount)
	router.GET("/api/users/:id", server.getUserAccountByID)
	router.GET("/api/users", server.listUsers)
	router.PUT("/api/users/:id", server.updateUserAccount)
	router.DELETE("/api/users/:id", server.deleteUserAccount)

	// accounts routes
	router.POST("/api/accounts", server.createAccount)
	router.GET("/api/accounts/:id", server.getAccountByID)
	router.GET("/api/accounts/users/:id", server.getAccountByUserID)
	router.GET("/api/accounts", server.listAccounts)
	router.PUT("/api/accounts/:id", server.updateAccount)

	server.router = router

	return server
}

// StartServer starts the server on a specific address
func (srv *Server) StartServer(addr string) error {
	return srv.router.Run(addr)
}

// errorResponse provides an error message
func errorResponse(s string, err error) gin.H {
	return gin.H{"message": s, "error": err.Error()}
}
