package router

import (
	"net/http"
	"reflect"
	"strings"
)

func SetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		keys := reflect.ValueOf(req.Header).MapKeys()
		strkeys := make([]string, len(keys))
		for i := 0; i < len(keys); i++ {
			strkeys[i] = keys[i].String()
		}
		headres := strings.Join(strkeys, ",") + ", X-Naver-Client-Id,X-Naver-Client-Secret,X-TARGET-URL,Content-Type"

		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:5000")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Set("Access-Control-Allow-Methods", "*")
		res.Header().Set("Content-Type", "application/json") //res type json set
		res.Header().Set("Access-Control-Allow-Headers", headres)
		next.ServeHTTP(res, req)
	})
}
