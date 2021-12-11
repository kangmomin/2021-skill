# project 2021 skill contest

+ 포스트 정보 가져오기
    ```
    method: GET
    link: /post
    query: page int

    ex) /post?page=1
    ```
    
+ 각 게시글의 정보 가져오기
    ```
    method: GET
    link: /post/:id
    params: id int
    
    ex) /post/12
    ```
    
+ 계정 생성
    ```
    method: POST
    link: /sign-up
    body: {
        name      string
	    id        string
	    password  string
	    studentId string
    }
    
    ex) /sign-up
        body {
            name:      "example",
            id:        "hello",
            password:  "password",
            studentId: "10101",
        }
    ```