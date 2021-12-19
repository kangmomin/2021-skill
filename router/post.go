package router

import (
	"2021skill/account"
	"2021skill/conn"
	"2021skill/structure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//for write post
type writePostBody struct {
	Description string `json:"description"`
	Tocken      string `json:"tocken"`
	Title       string `json:"title"`
}

//for get post
type resJson struct {
	Data     []structure.DB `json:"data"`
	LastPage int            `json:"lastPage"`
	NowPage  int            `json:"nowPage"`
}

//for delete post and get each post
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
		var _err bool
		userId, _err = account.DecodeTocken(body.Tocken)

		if _err {
			res.WriteHeader(400)
			fmt.Fprint(res, "error during decode tocken")
			return
		}
	}

	//get post's info
	err := db.QueryRow("SELECT * FROM post WHERE id=?", postId).Scan(&post.Id, &post.Title, &post.Description, &post.Good, &post.Bad, &post.ReplyCount, &post.View, &post.Time, &post.OwnerId)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "error during get post's info")
		return
	}

	//check is owner
	post.IsOwner = false
	if userId := string(userId); userId == strconv.Itoa(post.OwnerId) {
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

func DeletePost(res http.ResponseWriter, req *http.Request) {
	postId := req.URL.Query()["id"]
	if len(postId) < 1 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "not enough querystring")
		return
	}

	var body tocken
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

	var postOwnerId int
	db.QueryRow("SELECT ownerId FROM post WHERE id=?", postId[0]).Scan(&postOwnerId)

	if userId, _ := strconv.Atoi(string(userId)); postOwnerId != userId {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "not a user's post")
		return
	}

	_, err := db.Exec("DELETE FROM post WHERE id=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	_, err = db.Exec("DELETE FROM reply WHERE postId=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	_, err = db.Exec("DELETE FROM good WHERE postId=?", postId[0])
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during deleting")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, "success")
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
	var posts resJson //result 변수
	db := conn.DB

	pageString := getPage(req)
	page, err := strconv.Atoi(pageString)

	//query pages info set
	posts.NowPage = page
	page = (page - 1) * 20 //each page's count
	eachPostConunt := 20

	//만약 queryString에 page값이 없으면 page를 0으로 셋
	if err != nil {
		page = 0
	}

	sort := req.URL.Query()["sort"]
	var sortType string

	if len(sort) < 1 || len(sort[0]) < 1 { //만약 sort type이 없으면 기본으로 넘기고 id값이면 오름차, 그외의 값은 내림차순 정렬
		sortType = "id DESC"
	} else {
		sortType = sort[0] + " DESC"
	}

	keyWord := ""
	if len(req.URL.Query()["search"]) > 0 {
		keyWord = req.URL.Query()["search"][0]
	}

	//30개의 값을 꺼내옴, colum에 맞는 정렬, 서칭
	query := "SELECT * FROM `post` WHERE title LIKE '%" + keyWord + "%' ORDER BY " + sortType + " LIMIT ?, ?;"
	post, err := db.Query(query, strconv.Itoa(page), strconv.Itoa(eachPostConunt))

	//조건에 맞는 패이지의 수
	var allPosts, lastPage int
	db.QueryRow("SELECT COUNT(*) FROM post WHERE title LIKE '%" + keyWord + "%';").Scan(&allPosts)

	lastPage = allPosts / eachPostConunt
	if (allPosts % eachPostConunt) != 0 {
		lastPage++
	}
	posts.LastPage = lastPage

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
		fmt.Println("error during json marshaling")
		fmt.Println(err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during json marshaling")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, string(postJson))
}

func WritePost(res http.ResponseWriter, req *http.Request) {
	//get body data
	var body writePostBody
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
	}

	//if body don't have tocken == user wasn't login
	if len(body.Tocken) < 1 {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "need Login")
		return
	}

	//if description or title is null
	if len(body.Description) < 1 || len(body.Title) < 1 || len(body.Title) > 30 {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, "some value is null")
		return
	}

	userId, _err := account.DecodeTocken(body.Tocken)

	if _err {
		res.WriteHeader(400)
		fmt.Fprint(res, "error during decode tocken")
		return
	}

	//insert post
	db := conn.DB
	body.Description = strings.ReplaceAll(body.Description, "script", "")
	_, err = db.Exec("INSERT INTO post (title, description, ownerId) VALUES (?, ?, ?)", body.Title, body.Description, userId)

	if err != nil {
		fmt.Printf("error during insert post \n %s", err)
		res.WriteHeader(400)
		fmt.Fprint(res, "error during insert post")
		return
	}

	res.WriteHeader(200)
	fmt.Fprint(res, "success")
}
