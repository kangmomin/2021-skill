package account

import (
	"2021skill/conn"
	"2021skill/logger"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetAccountList(res http.ResponseWriter, req *http.Request) {
	var body body
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		logger.ErrLogger().Fatalln(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "not allowed method")
		return
	}
	userId, _err := DecodeTocken(body.Tocken)
	if _err {
		logger.ErrLogger().Fatalln(_err)
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
	rows, err := db.Query("SELECT id, name, accountId, studentId, authImg FROM account WHERE accept=0")
	if err != nil {
		logger.ErrLogger().Fatalln(err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get data")
		return
	}

	var resValue []structure.Account
	for rows.Next() {
		var row structure.Account
		rows.Scan(&row.Id, &row.Name, &row.AccountId, &row.StudentId, &row.AuthImg)
		resValue = append(resValue, row)
	}

	resJson, _ := json.Marshal(resValue)
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, string(resJson))
}
