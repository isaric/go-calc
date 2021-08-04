package main

import (
	"go-calc/calc"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/{operation}", calc.Handle).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
