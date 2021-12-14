package router

import (
	"2021skill/conn"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lemon-mint/vbox"
)

type addGoodBody struct {
	Tocken string `json:"tocken"`
}

func AddGood(res http.ResponseWriter, req *http.Request) {
	var body addGoodBody
	json.NewDecoder(req.Body).Decode(&body)

	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need login")
		return
	}

	//decode tocken
	keyFile, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(keyFile)

	auth := vbox.NewBlackBox(key)                      //aes key
	decodeTocken, err := hex.DecodeString(body.Tocken) //tocken first decode

	userId, _error := auth.Open(decodeTocken)
	if !_error || err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "tocken was brocken")
		return
	}

	//get post id
	postId := mux.Vars(req)["id"]
	if postId == "" {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "id is null")
		return
	}

	//insert good info
	db := conn.DB
	_, err = db.Exec("INSERT INTO good (postId, userId) VALUES (?, ?)", postId, userId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during inerting")
		return
	}

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, "success")
}
