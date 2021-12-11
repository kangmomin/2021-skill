package router

import (
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type resJson struct {
	Data     []structure.DB `json:"data"`
	LastPage int            `json:"lastPage"`
	NowPage  int            `json:"nowPage"`
}

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

	//30개의 값을 꺼내옴
	query := "SELECT * FROM post LIMIT " + strconv.Itoa(page) + ", " + strconv.Itoa(eachPostConunt) + ";"
	post, err := db.Query(query)

	var posts resJson //result 변수

	//조건에 맞는 패이지의 수
	var allPosts, lastPage int
	db.QueryRow("SELECT COUNT(*) FROM post LIMIT " + strconv.Itoa(page) + ", " + strconv.Itoa(eachPostConunt) + ";").Scan(&allPosts)

	lastPage = allPosts / eachPostConunt
	if (allPosts % eachPostConunt) != 0 {
		lastPage++
	}
	posts.LastPage = lastPage
	posts.NowPage = page + 1 //query문을 위해 -1을 했지만 현재 패이지를 표시하기위해 +1

	if err != nil {
		fmt.Print(err)
		res.WriteHeader(500)
		fmt.Fprint(res, "error during get from db")
		return
	}

	//append dataes to posts
	for post.Next() {
		var row structure.DB
		err := post.Scan(&row.Id, &row.Title, &row.Description, &row.Good, &row.Bad, &row.ReplyCount, &row.View, &row.Time, &row.OwnerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		row.Created = row.Time.Format("06-01-02 15:04")
		posts.Data = append(posts.Data, row)
	}

	//value가 없으면 not found return
	if len(posts.Data) < 1 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "there are no posts")
		return
	}

	res.WriteHeader(200)
	postJson, err := json.Marshal(posts)
	if err != nil {
		panic(err)
	}

	res.WriteHeader(200)
	fmt.Fprint(res, string(postJson))
}
