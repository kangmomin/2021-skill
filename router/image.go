package router

import (
	"2021skill/account"
	"2021skill/conn"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type getAuthImgBody struct {
	ImgPath string
	Tocken  string
}

type resValue struct {
	ImgPath string `json:"imgPath"`
	Message string `json:"message"`
	Err     bool   `json:"error"`
}

func UploadImg(res http.ResponseWriter, req *http.Request) {
	//폼 파일 포맷의 image를 가져옴
	//header에는 img의 정보 img는 그냥 img만
	img, imgHeader, err := req.FormFile("image")
	var resValue resValue

	if err != nil {
		fmt.Println("err in uploadImg.go 25")

		resValue.Err = true
		resValue.Message = "error during get image data"

		fmt.Println(err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during get image data")
		return
	}

	//이미지 인코딩
	encodedImg, _ := ioutil.ReadAll(img)

	//이미지가 들어갈 위치와 이름 설정
	imgName := time.Now().String() + imgHeader.Filename
	imgName = strings.ReplaceAll(imgName, " ", "") //img이름 설정시 생기는 공백 제거
	imgName = strings.ReplaceAll(imgName, ":", "") //img이름 설정시 생기는 공백 제거

	//파일 업로드
	err = ioutil.WriteFile("public/"+imgName, encodedImg, 0644)

	if err != nil {
		fmt.Println(err)
		resValue.Message = "cannot uplaod the img"
		resValue.Err = true

		res.WriteHeader(400)
		fmt.Fprint(res, resValue)
		return
	}

	resValue.ImgPath = imgName
	resValue.Message = "success"
	resValue.Err = false

	resJson, _ := json.Marshal(resValue)

	fmt.Fprint(res, string(resJson))
}

func GetAuthImg(res http.ResponseWriter, req *http.Request) {
	var body getAuthImgBody

	json.NewDecoder(req.Body).Decode(&body)

	userId, _err := account.DecodeTocken(body.Tocken)

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

	img, err := os.Open("authImg/" + userAuthImg)

	if err != nil {
		fmt.Printf("error getAuthImg.go - 52 \n %s \n", err)
		fmt.Fprint(res, "error open img")
		return
	}

	res.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(res, img)
}

func UploadAuthImg(res http.ResponseWriter, req *http.Request) {
	//폼 파일 포맷의 image를 가져옴
	//header에는 img의 정보 img는 그냥 img만
	img, imgHeader, err := req.FormFile("image")
	var resValue resValue

	if err != nil {
		fmt.Println("err in uploadImg.go 25")

		resValue.Err = true
		resValue.Message = "error during get image data"

		fmt.Println(err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during get image data")
		return
	}

	//이미지 인코딩
	encodedImg, _ := ioutil.ReadAll(img)

	//이미지가 들어갈 위치와 이름 설정
	imgName := time.Now().String() + imgHeader.Filename
	imgName = strings.ReplaceAll(imgName, " ", "") //img이름 설정시 생기는 공백 제거
	imgName = strings.ReplaceAll(imgName, ":", "") //img이름 설정시 생기는 공백 제거

	//파일 업로드
	err = ioutil.WriteFile("authImg/"+imgName, encodedImg, 0644)

	if err != nil {
		fmt.Println(err)
		resValue.Message = "cannot uplaod the img"
		resValue.Err = true

		res.WriteHeader(400)
		fmt.Fprint(res, resValue)
		return
	}

	resValue.ImgPath = imgName
	resValue.Message = "success"
	resValue.Err = false

	resJson, _ := json.Marshal(resValue)

	fmt.Fprint(res, string(resJson))
}
