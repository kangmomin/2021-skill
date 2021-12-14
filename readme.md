# project 2021 skill contest

+ 포스트 정보 가져오기
    ```
    method: GET
    link: /post
    query: page int, sort string
    
    return = errorMessage string || [{
        id          int
        title       string
        description string
        good        int
        bad         int
        replyCount  int
        view        int
        created     string
        ownerId     int
    }]

    ex) /post?page=1
        
        return = [
            아래 각 게시글 정보 가져오기 부분의 return
        ]
    ```
    
+ 각 게시글의 정보 가져오기
    ```
    method: GET
    link: /post/:id
    params: id int
    return = errorMessage string || {
        id          int
        title       string
        description string
        good        int
        bad         int
        replyCount  int
        view        int
        created     string
        ownerId     int
    }
    
    ex) /post/12
    
        return = {
            id: 1,
            title: "hello",
            description: "example",
            good: 0,
            bad: 0,
            replyCount: 0,
            view: 34,
            created: "21-12-10 08:10",
            ownerId: 1
        }
    ```
    
+ 로그인
    ```
    method: POST
    link: /login
    body: {
        accountId       string //계정 아이디
        accountPassword string //계정 비밀번호
    }
    return = {
        userId  int
        message string
        err     boolean
    }
    
    ex) /login
        body: {
            accountId:       "exampleId",
            accountPassword: "examplePassword"
        }

        return = {
            userId:  0,
            message: "id is wrong",
            err:     true
        }
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
    return: insertedId int || errorMsg string
    
    ex) /sign-up
        body {
            name:      "example",
            id:        "hello",
            password:  "password",
            studentId: "10101"
        }
        
        return = 3
    ```
    
+ 중복 확인
    * 아이디
        ```
        method: POST
        link: /overlap-check/id
        body: {
            id string
        }
        return: boolean
        
        ex) /overlap-check/id
            body: {
                id: "example"
            }
            
            return = false
        ```

    * 학번
        ```
        method: POST
        link: /overlap-check/student-id
        body: {
            studentId int
        }
        return: boolean
        
        ex) /overlap-check/student-id
            body: {
                id: 10101
            }
            
            return = false
        ```

+ 이미지 업로드
    ```
    method: POST
    link: /upload
    body: formData()
    return: {
        error   bool
        message string
        imgPath string
    }
    
    ex) /upload
        body: {
            image: formData()
        }
        
        return = {
            error: true,
            message: "error during upload",
            imgPath: ""
        }
    ```

+ 글 작성
    ```
    method: POST
    link: /post
    body: {
        title       string
	    description string
	    tocken      string
    }
    return: errorMsg string || "success"
    
    ex) /sign-up
        body {
            title:       "hello World"
    	    description: "hello World" 
    	    tocken:      tocken
        }
        
        return = "error" || "success"
    ```