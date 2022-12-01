package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/blessedmadukoma/trackit-chima/api"
	db "github.com/blessedmadukoma/trackit-chima/db/sqlc"
	"github.com/blessedmadukoma/trackit-chima/util"
	_ "github.com/lib/pq"
)

func main() {
	// load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&timezone=%s", config.DB_DIALECT, config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME, config.DB_TIMEZONE)

	// connect to database
	conn, err := sql.Open(config.DB_DRIVER, dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
