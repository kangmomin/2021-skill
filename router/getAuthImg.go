package router

import (
	"2021skill/conn"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/lemon-mint/vbox"
)

type body struct {
	ImgPath string
	Auth    string
}

func GetAuthImg(res http.ResponseWriter, req *http.Request) {
	var body body

	json.NewDecoder(req.Body).Decode(&body)

	reader, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(reader)
	auth := vbox.NewBlackBox(key)

	userIdByte, _ := auth.Open([]byte(body.Auth))
	userId, err := strconv.Atoi(hex.EncodeToString(userIdByte))

	if err != nil {
		fmt.Printf("error getAuthImg.go - 34 \n %s \n", err)
		fmt.Fprint(res, "error during decode auth key")
		return
	}

	db := conn.DB
	var userAuthImg string

	err = db.QueryRow("SELECT authImg FROM account WHERE id=?", userId).Scan(&userAuthImg)

	if err != nil {
		fmt.Printf("error getAuthImg.go - 45 \n %s \n", err)
		fmt.Fprint(res, "error during get user's auth img src")
		return
	}

	img, err := os.Open("public/" + userAuthImg)

	if err != nil {
		fmt.Printf("error getAuthImg.go - 52 \n %s \n", err)
		fmt.Fprint(res, "error open img")
		return
	}

	res.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(res, img)
}
