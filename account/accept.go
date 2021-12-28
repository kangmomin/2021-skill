package account

import (
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type body struct {
	Tocken string `json:"tocken"`
	UserId int    `json:"userId"`
}

func Accept(res http.ResponseWriter, req *http.Request) {
	var body body
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "not allowed method")
		return
	}
	userId, _err := DecodeTocken(body.Tocken)
	if _err {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "not allowed method")
		return
	}
	if userId, _ := strconv.Atoi(string(userId)); userId != 33 {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "관리자가 아닙니다.")
		return
	}

	db := conn.DB

	_, err = db.Exec("UPDATE account SET accept=1 WHERE id=?", body.UserId)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during update")
		return
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, "accepted")
}
