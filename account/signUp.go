package account

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type post struct {
	Name      string `json:"name"`
	AccountId string `json:"id"`
	Password  string `json:"password"`
	StudentId string `json:"studentId"`
}

func SignUp(res http.ResponseWriter, req *http.Request) {
	db := conn.DB

	//get body
	body := post{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "error in body data")
		return
	}

	//빈칸 체크
	if len(body.AccountId) < 1 || len(body.Name) < 1 || len(body.Password) < 1 || len(body.StudentId) < 1 {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "empty value")
		return
	}

	//중복 체크를 위한 account db 가져오기
	rows, err := db.Query("SELECT accountId, name, studentId FROM account")

	if err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, "err during get data from db")
		fmt.Print(err)
		return
	}

	//중복 체크
	for rows.Next() {
		row := structure.Account{}
		rows.Scan(&row.Name, &row.AccountId, &row.StudentId)

		if row.Name == body.Name || row.AccountId == body.AccountId || strconv.Itoa(row.StudentId) == body.StudentId {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "something is overlaped")
			return
		}
	}

	//db에 계정 추가
	userInfo, err := db.Exec("INSERT INTO account (name, accountId, accountPassword, studentId) VALUES (?, ?, ?, ?)",
		body.Name, body.AccountId, body.Password, body.StudentId)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during inserting")
		fmt.Print(err)
		return
	}

	//추가된 계정 id값
	insertedID, _ := userInfo.LastInsertId()

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, strconv.FormatInt(insertedID, 10))
}
