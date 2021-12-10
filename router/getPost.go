package router

import (
	"2021skill/conn"
	"2021skill/structure"
	"fmt"
	"net/http"
	"strconv"
)

func getPage(req *http.Request) (page string) {
	//get queryString
	queryString := req.URL.Query()

	if len(queryString["page"]) < 1 {
		return "1"
	}

	return queryString["page"][0]
}

func GetPost(res http.ResponseWriter, req *http.Request) {
	db := conn.DB

	pageString := getPage(req)
	page, err := strconv.Atoi(pageString)

	//query pages info set
	page -= 1 //each page's count
	eachPostConunt := 30

	//만약 queryString에 page값이 없으면 page를 0으로 셋
	if err != nil {
		page = 0
	}

	query := "SELECT * FROM post LIMIT " + strconv.Itoa(page) + ", " + strconv.Itoa(eachPostConunt) + ";"
	post, err := db.Query(query)

	var posts []structure.DB

	if err != nil {
		fmt.Print(err)
		res.WriteHeader(500)
		fmt.Fprint(res, "error during get from db")
		return
	}

	//append dataes to posts
	for post.Next() {
		var row structure.DB
		err := post.Scan(&row.Id, &row.Title, &row.Description, &row.Good, &row.Bad, &row.ReplyCount, &row.View, &row.Created)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, row)
	}

	//value가 없으면 not found return
	if len(posts) < 1 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "there are no posts")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, posts)
}
