package main

import (
	"2021skill/account"
	"2021skill/router"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	app := mux.NewRouter()

	app.Use(router.SetHeader)

	app.HandleFunc("/post", router.GetPost).Methods("GET")
	app.HandleFunc("/post/{id}", router.GetEachPost).Methods("GET")

	app.HandleFunc("/sign-up", account.SignUp).Methods("POST")

	http.ListenAndServe(":3101", app)
}
