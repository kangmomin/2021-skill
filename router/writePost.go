package router

import (
	"2021skill/conn"
	"encoding/hex"
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
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
	}

	//if body don't have tocken == user wasn't login
	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need Login")
		return
	}

	//if description or title is null
	if len(body.Description) < 1 || len(body.Title) < 1 || len(body.Title) > 30 {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "some value is null")
		return
	}

	//get key form config folder
	file, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(file)

	//decode userInfo
	auth := vbox.NewBlackBox(key)
	decodedTocken, _ := hex.DecodeString(body.Tocken)
	userId, boolean := auth.Open(decodedTocken)

	if err != nil || !boolean {
		fmt.Println(req.Cookies())
		fmt.Printf("error during decode tocken \n %s", err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
		return
	}

	//insert post
	db := conn.DB
	_, err = db.Exec("INSERT INTO post (title, description, ownerId) VALUES (?, ?, ?)", body.Title, body.Description, userId)

	if err != nil {
		fmt.Printf("error during insert post \n %s", err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during insert post")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, "success")
}
