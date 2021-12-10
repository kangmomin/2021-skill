package router

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetEachPost(res http.ResponseWriter, req *http.Request) {
	db := conn.DB
	var post structure.DB

	//get params
	postId := mux.Vars(req)["id"]

	//get post's info
	err := db.QueryRow("SELECT * FROM post WHERE id=?", postId).Scan(&post.Id, &post.Title, &post.Description, &post.Good, &post.Bad, &post.ReplyCount, &post.View, &post.Time)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error")
		return
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

	res.Header().Set("Content-Type", "application/json") //res type json set
	res.WriteHeader(200)
	fmt.Fprint(res, string(byteDB))
}
