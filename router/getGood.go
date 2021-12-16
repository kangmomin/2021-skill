package router

import (
	"2021skill/account"
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	userId, _err := account.DecodeTocken(body.Tocken)

	if _err {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
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
	err := db.QueryRow("SELECT COUNT(*) FROM good WHERE userId=? AND postId=?", userId, postId).Scan(&resValue.IsGood)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "err during get data")
		return
	}

	resJson, _ := json.Marshal(resValue)
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, string(resJson))
}
