package router

import (
	"2021skill/conn"
	"2021skill/structure"
	"fmt"
	"net/http"
	"strconv"
)

func GetPost(res http.ResponseWriter, req *http.Request) {
	db := conn.DB

	//get queryString
	queryString := req.URL.Query()
	pageString := queryString["page"][0]
	page, err := strconv.Atoi(pageString)

	//query pages info set
	page -= 1 //each page's count
	eachPostConunt := 30

	if err != nil {
		panic(err.Error())
	}

	post, err := db.Query("SELECT * FROM post LIMIT ? ?", page, eachPostConunt)

	var posts []structure.DB

	if err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, "error during get from db")
	}

	//append dataes to posts
	for post.Next() {
		var row structure.DB
		post.Scan(&row.Id, &row.Title, &row.Description, &row.Good,
			&row.Bad, &row.ReplyCount, &row.View, &row.Created)
		posts = append(posts, row)
	}

	if len(posts) < 1 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "there are no posts")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, posts)
}
