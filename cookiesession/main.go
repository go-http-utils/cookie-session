package main

import (
	"net/http"
	"net/http/httptest"

	cookiesession "github.com/go-http-utils/cookie-session"
)

func main() {
	req, _ := http.NewRequest("GET", "/health-check", nil)
	store := cookiesession.NewCookieStore(nil)

	cookiekey := "teambition"

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, cookiekey)
		session.Values["name"] = "mushroom"
		session.Values[66] = 99
		session.Save(r, w)
	})
	handler.ServeHTTP(recorder, req)

	cookies, _ := getCookie(cookiekey, recorder)

	//======reuse=====
	store = cookiesession.NewCookieStore(nil)
	req.AddCookie(cookies)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store.Get(r, cookiekey)
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
