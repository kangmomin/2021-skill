package router

import (
	"2021skill/account"
	"2021skill/logger"
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
	ImgPath string `json:"imgPath"`
	Tocken  string `json:"tocken"`
}

type resValue struct {
	ImgPath string `json:"imgPath"`
	Message string `json:"message"`
	Err     bool   `json:"error"`
}

func ImageUploader(res http.ResponseWriter, req *http.Request) {
	//폼 파일 포맷의 image를 가져옴
	//header에는 img의 정보 img는 그냥 img만
	img, imgHeader, err := req.FormFile("image")
	var resValue resValue

	if err != nil {
		resValue.Err = true
		resValue.Message = "error during get image data"

		logger.ErrLogger().Fatalln(err)
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
	fileType := "public/"

	if req.URL.Path == "/auth-image" {
		fileType = "auth/"
	}

	err = ioutil.WriteFile(fileType+imgName, encodedImg, 0644)

	if err != nil {
		resValue.Message = "cannot uplaod the img"
		resValue.Err = true

		logger.ErrLogger().Fatalln(err)
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
		logger.ErrLogger().Fatalln(_err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
		return
	}

	if string(userId) != "33" {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "not admin")
		return
	}

	img, err := os.Open("authImg/" + body.ImgPath)

	if err != nil {
		logger.ErrLogger().Fatalln(err)
		fmt.Fprint(res, "error open img")
		return
	}

	res.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(res, img)
}
