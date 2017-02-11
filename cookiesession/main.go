package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-http-utils/cookie-session"
)

func main() {

	sessionkey := "sessionid"

	store := sessions.NewCookieStore([]string{})

	req, _ := http.NewRequest("GET", "/", nil)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(sessionkey, w, r)
		if val, ok := session.Values["name"]; ok {
			println(val)
		} else {
			session.Values["name"] = "mushroom"
		}
		session.Save()
	})
	handler.ServeHTTP(recorder, req)

	//======reuse=====
	store = sessions.NewCookieStore([]string{})
	cookies, _ := getCookie(sessionkey, recorder)
	req.AddCookie(cookies)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(sessionkey, w, r)

		println(session.Values["name"].(string))
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
