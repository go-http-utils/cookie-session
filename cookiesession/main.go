package main

import (
	"net/http"
	"net/http/httptest"

	cookiesession "github.com/go-http-utils/cookie-session"
)

func main() {
	req, _ := http.NewRequest("GET", "/health-check", nil)

	cookiekey := "teambition"

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := cookiesession.New(w, r)
		session, _ := store.Get(cookiekey)
		session.Values["name"] = "mushroom"
		session.Values[66] = 99
		session.Save()
	})
	handler.ServeHTTP(recorder, req)

	cookies, _ := getCookie(cookiekey, recorder)

	//======reuse=====

	req.AddCookie(cookies)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := cookiesession.New(w, r)
		session, _ := store.Get(cookiekey)
		println(session.Values["name"].(string))
		println(session.Values[66].(int))
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
