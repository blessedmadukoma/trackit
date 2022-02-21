package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)

	fmt.Println("Server starting on port 8080!")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
