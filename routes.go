package main

import (
	"log"

	"github.com/gorilla/mux"
)

// AddApproutes will add the routes for the application
func AddApproutes(route *mux.Router) {

	log.Println("Loadeding Routes...")

	route.HandleFunc("/signup", SignUpUser).Methods("POST")

	route.HandleFunc("/resetPassword", ResetPassword).Methods("POST")

	log.Println("Routes are Loaded.")
}
