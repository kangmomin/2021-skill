package account

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/lemon-mint/vbox"
	"golang.org/x/crypto/argon2"
)

type resStruct struct {
	UserId  int    `json:"userId"`
	Message string `json:"message"`
	Token   string `json:"tocken"`
	Err     bool   `json:"err"`
}

func Login(res http.ResponseWriter, req *http.Request) {
	var resValue resStruct
	var body structure.Account
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
		fmt.Println(req.Body)
		res.WriteHeader(http.StatusBadRequest)
		resValue.Err = true
		resValue.Message = "error in body data"
		resJson, _ := json.Marshal(resValue)
		fmt.Fprint(res, string(resJson))
		return
	}

	db := conn.DB
	var accountPassword, encryptedPassword, random string //인증 관련 변수들
	var userId int

	err = db.QueryRow("SELECT random, accountPassword, id FROM account WHERE accountId=?", body.AccountId).Scan(&random, &accountPassword, &userId)

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(400)
		resValue.Err = true
		resValue.Message = "id is wrong"
		resJson, _ := json.Marshal(resValue)
		fmt.Fprint(res, string(resJson))
		return
	}

	resValue.UserId = userId

	//body로 넘어온 password 암호화
	salt, _ := hex.DecodeString(random) //db에 저장하기 위해 encode했던 string을 byte로 decode
	byteEncryptPwd := argon2.IDKey([]byte(body.Password), []byte(salt), structure.EncryptConfig.Iterations,
		structure.EncryptConfig.Memory, structure.EncryptConfig.Parallelism, structure.EncryptConfig.KeyLength)
	encryptedPassword = hex.EncodeToString(byteEncryptPwd)

	//암호화 되어있는 두 값을 비교
	if encryptedPassword != accountPassword {
		res.WriteHeader(400)
		resValue.Err = true
		resValue.Message = "password is wrong"
		resJson, _ := json.Marshal(resValue)
		fmt.Fprint(res, string(resJson))
		return
	}

	//현재 key
	keyFile, err := os.Open("config/accessKey.txt")

	if err != nil {
		resValue.Err = true
		res.WriteHeader(400)
		fmt.Println("error during get accessKey config file")
		fmt.Println(err)
	}
	key, _ := ioutil.ReadFile(keyFile.Name())

	aesKey := vbox.NewBlackBox(key) //aes 암호화 key

	//data encrypt
	auth := aesKey.Seal([]byte(strconv.Itoa(userId)))
	resValue.Message = "login success"
	resValue.Token = hex.EncodeToString(auth)
	resJson, _ := json.Marshal(resValue) //object to json

	resValue.Err = false
	res.WriteHeader(200)
	fmt.Fprint(res, string(resJson))
}
