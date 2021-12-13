package account

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"golang.org/x/crypto/argon2"
)

type post struct {
	Name      string `json:"name"`
	AccountId string `json:"id"`
	Password  string `json:"password"`
	StudentId string `json:"studentId"`
	AuthImg   string `json:"authImg"`
}

func SignUp(res http.ResponseWriter, req *http.Request) {
	db := conn.DB

	//get body
	body := post{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "error in body data")
		return
	}

	//빈칸 체크
	if len(body.AccountId) < 1 || len(body.Name) < 1 || len(body.Password) < 1 || len(body.StudentId) < 1 || len(body.AuthImg) < 1 {
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
		rows.Scan(&row.AccountId, &row.Name, &row.StudentId)

		if row.Name == body.Name || row.AccountId == body.AccountId || strconv.Itoa(row.StudentId) == body.StudentId {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "something is overlaped")
			return
		}
	}

	//password encrypting
	randomKey, _ := generateRandomBytes(structure.EncryptConfig.SaltLength) //salt
	encryptedPwd := argon2.IDKey([]byte(body.Password), randomKey, structure.EncryptConfig.Iterations,
		structure.EncryptConfig.Memory, structure.EncryptConfig.Parallelism, structure.EncryptConfig.KeyLength)
	body.Password = hex.EncodeToString(encryptedPwd)

	//db에 계정 추가
	userInfo, err := db.Exec("INSERT INTO account (name, accountId, accountPassword, studentId, random, authImg) VALUES (?, ?, ?, ?, ?, ?)",
		body.Name, body.AccountId, body.Password, body.StudentId, hex.EncodeToString(randomKey), body.AuthImg)
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

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
