package account

import (
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"net/http"
)

//body값을 받기위한 structure
type Id struct {
	Id string
}

type StuendtId struct {
	StuendtId string
}

func IdOverLap(res http.ResponseWriter, req *http.Request) {
	var body Id
	var overLapId int
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		panic(err)
	}
	db := conn.DB

	err = db.QueryRow("SELECT id FROM account WHERE accountId=?", body.Id).Scan(&overLapId)

	if errString := err.Error(); err != nil && errString != "sql: no rows in result set" { //가져온 값이 없더라도 넘어감.
		panic(err)
	}

	res.WriteHeader(200)
	if overLapId == 0 { //overlap값이 없을때 false를 리턴
		fmt.Fprint(res, false)
		return
	}
	fmt.Fprint(res, true)
}

//사실상 같은 코드이지만 생각하기가 귀찮았음...
func StuendtIdOverLap(res http.ResponseWriter, req *http.Request) {
	var body StuendtId
	var overLapStuendtId int
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		panic(err)
	}
	db := conn.DB

	//db에서 입력된 학번과 일치하는 값이 있는지 확인
	err = db.QueryRow("SELECT id FROM account WHERE stuendtId=?", body.StuendtId).Scan(&overLapStuendtId)

	if errString := err.Error(); err != nil && errString != "sql: no rows in result set" { //가져온 값이 없더라도 넘어감.
		panic(err)
	}

	res.WriteHeader(200)
	if overLapStuendtId == 0 { //overlap값이 없을때 false를 리턴
		fmt.Fprint(res, false)
		return
	}
	fmt.Fprint(res, true)
}
