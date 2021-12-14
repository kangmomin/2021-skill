package router

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetReply(res http.ResponseWriter, req *http.Request) {
	postId := mux.Vars(req)["id"]
	if len(postId) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "need post id")
		return
	}

	db := conn.DB

	reply, err := db.Query("SELECT * FROM reply WHERE postId=?", postId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get reply data")
		return
	}
	var replies []structure.Reply

	for reply.Next() {
		var row structure.Reply

		reply.Scan(&row.Id, &row.OwnerId, &row.Description, &row.PostId, &row.Time)
		row.Created = row.Time.Format("06-01-02 15:04")
		replies = append(replies, row)
	}

	resJson, _ := json.Marshal(replies)

	res.WriteHeader(200)
	fmt.Fprint(res, string(resJson))
}
