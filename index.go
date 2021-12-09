package main

import (
	"2021skill/router"
	"net/http"
)

func main() {
	app := http.NewServeMux()

	app.HandleFunc("/post", router.GetPost)

	http.ListenAndServe(":3000", app)
}
