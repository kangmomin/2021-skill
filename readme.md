# project 2021 skill contest

+ 포스트 정보 가져오기
    ```
    link = /post
    query = page int

    ex) /post?page=1
    ```
    
+ 계정 생성
    ```
    link = /sign-up
    body = {
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