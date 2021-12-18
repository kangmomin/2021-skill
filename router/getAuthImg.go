package router

import (
	"2021skill/account"
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type body struct {
	ImgPath string
	Auth    string
}

func GetAuthImg(res http.ResponseWriter, req *http.Request) {
	var body body

	json.NewDecoder(req.Body).Decode(&body)

	userId, _err := account.DecodeTocken(body.Auth)

	if _err {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
		return
	}

	db := conn.DB
	var userAuthImg string

	err := db.QueryRow("SELECT authImg FROM account WHERE id=?", userId).Scan(&userAuthImg)

	if err != nil {
		fmt.Printf("error getAuthImg.go - 45 \n %s \n", err)
		fmt.Fprint(res, "error during get user's auth img src")
		return
	}

	img, err := os.Open("auth/" + userAuthImg)

	if err != nil {
		fmt.Printf("error getAuthImg.go - 52 \n %s \n", err)
		fmt.Fprint(res, "error open img")
		return
	}

	res.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(res, img)
}
