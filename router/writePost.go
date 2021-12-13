package router

import (
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lemon-mint/vbox"
)

type writePostBody struct {
	Description string `json:"description"`
	Tocken      string `json:"tocken"`
	Title       string `json:"title"`
}

func WritePost(res http.ResponseWriter, req *http.Request) {
	//get body data
	var body writePostBody
	json.NewDecoder(req.Body).Decode(&body)

	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need Login")
		return
	}

	if len(body.Description) < 1 || len(body.Title) < 1 {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "some value is null")
		return

	}

	//get key form config folder
	file, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(file)

	//decode userInfo
	auth := vbox.NewBlackBox(key)
	byteUserId, _ := auth.Open([]byte(body.Tocken))
	userId, _ := json.Marshal(byteUserId)

	db := conn.DB
	_, err := db.Exec("INSERT INFO post (title, description, ownerId) VALUES (?, ?, ?)", body.Title, body.Description, userId)

	if err != nil {
		fmt.Printf("error writePost.go 44 \n %s", err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during insert post")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, "success")
}
