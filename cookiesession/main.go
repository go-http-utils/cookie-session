package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-http-utils/cookie-session"
)

func main() {

	sessionkey := "sessionid"
	store := sessions.New([]string{"key"})

	req, _ := http.NewRequest("GET", "/", nil)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store.Init(w, r, false)
		session, _ := sessions.Get(sessionkey, store)
		session.Values["name"] = "mushroom"

		session.Save()
	})
	handler.ServeHTTP(recorder, req)

	//======reuse=====
	cookies, _ := getCookie(sessionkey, recorder)
	req.AddCookie(cookies)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store.Init(w, r, false)
		session, _ := sessions.Get(sessionkey, store)
		println(session.Values["name"].(string))
		println(session.Values["num"].(int))
	})
	handler.ServeHTTP(recorder, req)
}
func getCookie(name string, recorder *httptest.ResponseRecorder) (*http.Cookie, error) {
	var err error
	res := &http.Response{Header: http.Header{"Set-Cookie": recorder.HeaderMap["Set-Cookie"]}}
	for _, val := range res.Cookies() {
		if val.Name == name {
			return val, nil
		}
	}
	return nil, err
}
