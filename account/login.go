package account

import (
	"2021skill/conn"
	"2021skill/structure"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"

	session "github.com/go-session/session/v3"
)

type resStruct struct {
	userId int
}

func Login(res http.ResponseWriter, req *http.Request) {
	var body structure.Account
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error in body data")
		return
	}

	db := conn.DB
	var accountPassword, encryptedPassword, random string //인증 관련 변수들
	var userId int
	var resValue resStruct

	err = db.QueryRow("SELECT random, accountPassword, id FROM account WHERE accountId=?", body.AccountId).Scan(&random, &accountPassword, &userId)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(400)
		fmt.Fprint(res, "id is wrong")
		return
	}

	resValue.userId = userId

	//body로 넘어온 password 암호화
	salt, _ := hex.DecodeString(random) //db에 저장하기 위해 encode했던 string을 byte로 decode
	byteEncryptPwd := argon2.IDKey([]byte(body.Password), []byte(salt), structure.EncryptConfig.Iterations,
		structure.EncryptConfig.Memory, structure.EncryptConfig.Parallelism, structure.EncryptConfig.KeyLength)
	encryptedPassword = hex.EncodeToString(byteEncryptPwd)

	//암호화 되어있는 두 값을 비교
	if encryptedPassword != accountPassword {
		res.WriteHeader(400)
		fmt.Fprint(res, "password is wrong")
		return
	}

	//로그인 성공시 해당 계정의 session cookie생성
	store, err := session.Start(context.Background(), res, req)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(500)
		fmt.Fprint(res, "session error")
		return
	}

	store.Set("id", userId) //session에 유저의 id값을 넣음.

	res.WriteHeader(200)
	fmt.Fprint(res, resValue)
}
