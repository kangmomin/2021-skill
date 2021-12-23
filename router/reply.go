package router

import (
	"2021skill/account"
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//for write reply
type writeReplyBody struct {
	Tocken      string `json:"tocken"`
	Description string `json:"description"`
	PostId      string `json:"postId"`
}

func GetReply(res http.ResponseWriter, req *http.Request) {
	postId := mux.Vars(req)["id"]
	if len(postId) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "need post id")
		return
	}

	db := conn.DB

	reply, err := db.Query("SELECT * FROM reply WHERE postId=? AND refReplyId=0", postId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get reply data")
		return
	}

	refReply, err := db.Query("SELECT * FROM reply WHERE postId=? AND refReplyId!=0 ORDER BY id", postId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get ref_reply data")
		return
	}

	var replies structure.GetReply

	for reply.Next() {
		var row structure.Reply
		var refReplyId int

		reply.Scan(&row.Id, &row.OwnerId, &row.Description, &row.PostId, &refReplyId, &row.Time)
		row.Created = row.Time.Format("06-01-02 15:04")
		replies.Reply = append(replies.Reply, row)
	}

	for refReply.Next() {
		var row structure.RefReply

		refReply.Scan(&row.Id, &row.OwnerId, &row.Description, &row.PostId, &row.RefReplyId, &row.Time)
		row.Created = row.Time.Format("06-01-02 15:04")
		replies.RefReply = append(replies.RefReply, row)
	}

	resJson, _ := json.Marshal(replies)

	res.WriteHeader(200)
	fmt.Fprint(res, string(resJson))
}

func WriteReply(res http.ResponseWriter, req *http.Request) {
	var body writeReplyBody
	json.NewDecoder(req.Body).Decode(&body)

	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need login")
		return
	}

	userId, _err := account.DecodeTocken(body.Tocken)

	if _err {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
		return
	}

	db := conn.DB

	body.Description = strings.ReplaceAll(body.Description, "script", "")
	if len(body.Description) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "description is null")
		return
	}

	_, err := db.Exec("UPDATE post SET replyCount=replyCount+1 WHERE id=?", body.PostId)
	if err != nil {
		res.WriteHeader(400)
		fmt.Println("error during inserting")
		return
	}

	_, err = db.Exec("INSERT INTO reply (ownerId, description, postId) VALUES (?, ?, ?)", userId, body.Description, body.PostId)
	if err != nil {
		res.WriteHeader(400)
		fmt.Println("error during inserting")
		return
	}

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, "success")
}
