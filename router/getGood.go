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

type getGoodBody struct {
	Tocken string `json:"tocken"`
}

type getGoodRes struct {
	IsGood int `json:"isGood"`
}

func GetGood(res http.ResponseWriter, req *http.Request) {
	var body getGoodBody
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

	postId := "" //postId init
	postId = mux.Vars(req)["id"]
	if len(postId) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "need post id")
		return
	}

	db := conn.DB
	var resValue getGoodRes
	err = db.QueryRow("SELECT COUNT(*) FROM good WHERE userId=? AND postId=?", userId, postId).Scan(&resValue.IsGood)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "err during get data")
		return
	}

	resJson, _ := json.Marshal(resValue)
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, string(resJson))
}
