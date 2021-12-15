package router

import (
	"2021skill/conn"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/lemon-mint/vbox"
)

type deletePostBody struct {
	Tocken string `json:"tocken"`
}

func DeletePost(res http.ResponseWriter, req *http.Request) {
	postId := req.URL.Query()["id"]
	if len(postId) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "not enough querystring")
		return
	}

	var body deletePostBody
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

	var postOwnerId int
	db.QueryRow("SELECT ownerId FROM post WHERE id=?", postId[0]).Scan(&postOwnerId)

	if userId, _ := strconv.Atoi(string(userId)); postOwnerId != userId {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "not a user's post")
		return
	}

	_, err = db.Exec("DELETE FROM post WHERE id=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	_, err = db.Exec("DELETE FROM reply WHERE postId=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	_, err = db.Exec("DELETE FROM good WHERE postId=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, "success")
}
