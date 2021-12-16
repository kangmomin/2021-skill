package router

import (
	"2021skill/account"
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type deleteGoodBody struct {
	Tocken string `json:"tocken"`
}

func DeleteGood(res http.ResponseWriter, req *http.Request) {
	var body deleteGoodBody
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

	//get post id
	postId := mux.Vars(req)["id"]
	if postId == "" {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "id is null")
		return
	}

	//insert good info
	db := conn.DB
	_, err := db.Exec("DELETE FROM good WHERE postId=? AND userId=?", postId, userId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during inerting")
		return
	}

	_, err = db.Exec("UPDATE post SET good=good-1 WHERE id=?", postId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during inerting")
		return
	}

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, "success")
}
