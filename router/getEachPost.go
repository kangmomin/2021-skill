package router

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
	"time"

	"github.com/gorilla/mux"
	"github.com/lemon-mint/vbox"
)

type tocken struct {
	Tocken string `json:"tocken"`
}

func GetEachPost(res http.ResponseWriter, req *http.Request) {
	db := conn.DB
	var post structure.DB
	var body tocken
	var userId []byte

	//get params
	postId := mux.Vars(req)["id"]
	json.NewDecoder(req.Body).Decode(&body)

	//decode user info
	if len(body.Tocken) > 1 {
		file, _ := os.Open("config/accessKey.txt")
		key, _ := ioutil.ReadAll(file)

		auth := vbox.NewBlackBox(key)
		tocken, _ := hex.DecodeString(body.Tocken)
		userId, _ = auth.Open(tocken)
	}

	//get post's info
	err := db.QueryRow("SELECT * FROM post WHERE id=?", postId).Scan(&post.Id, &post.Title, &post.Description, &post.Good, &post.Bad, &post.ReplyCount, &post.View, &post.Time, &post.OwnerId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get post's info")
		return
	}

	//check is owner
	post.IsOwner = true
	if userId, _ := strconv.Atoi(hex.EncodeToString(userId)); userId == post.OwnerId {
		post.IsOwner = true
	}

	post.Created = post.Time.Format("06-01-02 15:04")
	byteDB, err := json.Marshal(&post)
	if err != nil {
		panic(err.Error())
	}

	//view count update
	_, err = db.Exec("UPDATE post SET view=? WHERE id=?", post.View+1, post.Id)

	if err != nil {
		fmt.Print(time.Now())
		fmt.Println(err)
	}

	res.WriteHeader(200)
	fmt.Fprint(res, string(byteDB))
}
