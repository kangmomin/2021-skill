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

	//post
	app.HandleFunc("/post", router.GetPost).Methods("GET")
	app.HandleFunc("/post/{id}", router.GetEachPost).Methods("POST")
	app.HandleFunc("/post", router.WritePost).Methods("POST")

	//account
	app.HandleFunc("/sign-up", account.SignUp).Methods("POST")
	app.HandleFunc("/login", account.Login).Methods("POST")

	//account overlap checker
	app.HandleFunc("/overlap-check/id", account.IdOverLap).Methods("POST")
	app.HandleFunc("/overlap-check/student-id", account.StuendtIdOverLap).Methods("POST")

	//image upload
	app.HandleFunc("/upload", router.UploadImg).Methods("POST")

	http.ListenAndServe(":3101", app)
}
