package account

import (
	"2021skill/conn"
	"2021skill/logger"
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

type response struct {
	IsOverlap bool
}

func IdOverLap(res http.ResponseWriter, req *http.Request) {
	var resVal response
	var body Id
	var overLapId int
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		panic(err)
	}
	db := conn.DB

	err = db.QueryRow("SELECT id FROM account WHERE accountId=?", body.Id).Scan(&overLapId)

	if err != nil { //가져온 값이 없더라도 넘어감.
		if errString := err.Error(); errString != "sql: no rows in result set" {
			logger.ErrLogger().Fatalln(err)
		}
	}

	res.WriteHeader(200)
	if overLapId == 0 { //overlap값이 없을때 false를 리턴
		resVal.IsOverlap = false
		resJSON, _ := json.Marshal(resVal)
		fmt.Fprint(res, string(resJSON))
		return
	}
	resVal.IsOverlap = true
	resJSON, _ := json.Marshal(resVal)
	fmt.Fprint(res, string(resJSON))
}

//사실상 같은 코드이지만 생각하기가 귀찮았음...
func StuendtIdOverLap(res http.ResponseWriter, req *http.Request) {
	var resVal response
	var body StuendtId
	var overLapStuendtId int
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		logger.ErrLogger().Fatalln(err)
	}

	db := conn.DB

	//db에서 입력된 학번과 일치하는 값이 있는지 확인
	err = db.QueryRow("SELECT id FROM account WHERE studentId=?", body.StuendtId).Scan(&overLapStuendtId)

	if err != nil { //가져온 값이 없더라도 넘어감.
		if errString := err.Error(); errString != "sql: no rows in result set" {
			logger.ErrLogger().Fatalln(err)
		}
	}

	res.WriteHeader(200)
	if overLapStuendtId == 0 { //overlap값이 없을때 false를 리턴
		resVal.IsOverlap = false
		resJSON, _ := json.Marshal(resVal)
		fmt.Fprint(res, string(resJSON))
		return
	}
	resVal.IsOverlap = true
	resJSON, _ := json.Marshal(resVal)
	fmt.Fprint(res, string(resJSON))
}
