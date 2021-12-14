package router

import (
	"2021skill/conn"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/lemon-mint/vbox"
)

type writeReplyBody struct {
	Tocken      string `json:"tocken"`
	Description string `json:"description"`
	PostId      int    `json:"postId"`
}

func WriteReply(res http.ResponseWriter, req *http.Request) {
	var body writeReplyBody
	json.NewDecoder(req.Body).Decode(&body)

	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need login")
		return
	}

	keyFile, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(keyFile)

	auth := vbox.NewBlackBox(key)
	decodeTocken, err := hex.DecodeString(body.Tocken)

	userId, _error := auth.Open(decodeTocken)
	if !_error || err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "tocken was brocken")
		return
	}

	db := conn.DB

	body.Description = strings.ReplaceAll(body.Description, "script", "")
	_, err = db.Exec("INSERT INTO reply (ownerId, description, postId) VALUES (?, ?, ?)", userId, body.Description, body.PostId)

	if err != nil {
		res.WriteHeader(400)
		fmt.Println("error during inserting")
		return
	}

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, "success")
}
