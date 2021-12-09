package main

import (
	"2021skill/router"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/post", router.GetPost)
}
