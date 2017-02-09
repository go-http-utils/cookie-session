package sessions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/cookie-session"
	"github.com/stretchr/testify/assert"
)

func TestGearsession(t *testing.T) {

	cookiekey := "teambition"
	cookieNewKey := "teambition-new"
	t.Run("gearsession use default options that should be", func(t *testing.T) {
		store := sessions.New()
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, false)
			session, _ := sessions.Get(cookiekey, store)
			session.Values["name"] = "mushroom"
			session.Values["num"] = 99
			session.Save()

		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		cookies, err := getCookie(cookiekey, recorder)
		assert.Nil(err)
		assert.NotNil(cookies.Value)
		t.Log(cookies.Value)
		req, err = http.NewRequest("GET", "/", nil)

		req.AddCookie(cookies)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, false)
			session, _ := sessions.Get(cookiekey, store)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(float64(99), session.Values["num"])
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		req, err = http.NewRequest("GET", "/", nil)
		req.AddCookie(cookies)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, false)
			session, _ := sessions.Get(cookiekey, store)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(float64(99), session.Values["num"])
		})
		handler.ServeHTTP(recorder, req)
	})
	store := sessions.New([]string{"key"})
	t.Run("gearsession with New session that should be", func(t *testing.T) {
		assert := assert.New(t)

		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, true)
			session, _ := sessions.Get(cookiekey, store)

			session.Values["name"] = "mushroom"
			session.Values["num"] = 99
			session.Save()

			session = session.New(cookieNewKey)
			session.Values["name"] = "teambition-n"
			session.Values["num"] = 100
			session.Save()
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====

		req, _ = http.NewRequest("GET", "/", nil)
		cookies, _ := getCookie(cookiekey, recorder)
		req.AddCookie(cookies)
		cookies, _ = getCookie(cookiekey+".sig", recorder)
		req.AddCookie(cookies)

		cookies, _ = getCookie(cookieNewKey, recorder)
		req.AddCookie(cookies)
		cookies, _ = getCookie(cookieNewKey+".sig", recorder)
		req.AddCookie(cookies)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, true)
			session, _ := sessions.Get(cookiekey, store)

			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(float64(99), session.Values["num"])

			session = session.Get(cookieNewKey)
			assert.Equal("teambition-n", session.Values["name"])
			assert.Equal(float64(100), session.Values["num"])

		})
		handler.ServeHTTP(recorder, req)

	})
	t.Run("gearsession with Name() and Store()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, false, &sessions.Options{
				Path:     "xxx.com",
				HTTPOnly: true,
				MaxAge:   64,
				Domain:   "ttt.com",
				Secure:   true,
			})
			session, _ := sessions.Get(cookiekey, store)

			session.Values["name"] = "mushroom"
			session.Values["num"] = 99
			session.Save()

			assert.Equal(cookiekey, session.Name())
			assert.NotNil(session.Store())

		})
		handler.ServeHTTP(recorder, req)
		cookies, _ := getCookie(cookiekey, recorder)
		assert.Equal("ttt.com", cookies.Domain)
		assert.Equal("xxx.com", cookies.Path)
		assert.Equal(true, cookies.HttpOnly)
		assert.Equal(64, cookies.MaxAge)
		assert.Equal(true, cookies.Secure)

	})
	t.Run("gearsession with Encode() and Decode()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		dt := make(map[string]interface{})
		err := sessions.Decode("xx", dt)
		assert.NotNil(err)

	})

}
func TestSessionCompatible(t *testing.T) {
	store := sessions.New([]string{"tb-accounts"})
	cookiekey := "TEAMBITION_SESSIONID"

	t.Run("gearsession should be compatible with old session component that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Cookie", "TEAMBITION_SESSIONID=eyJhdXRoVXBkYXRlZCI6MTQ4NTE1ODg3NDgxMywibmV4dFVybCI6Imh0dHA6Ly9wcm9qZWN0LmNpL3Byb2plY3RzIiwidHMiOjE0ODY2MDkzNTA5NjAsInVpZCI6IjU1YzE3MTBkZjk2YmJlODQ3NjgzMjUyYSIsInVzZXIiOnsiYXZhdGFyVXJsIjoiaHR0cDovL3N0cmlrZXIucHJvamVjdC5jaS90aHVtYm5haWwvMDEwa2UyZTMzODQ3ZjQzNzhlY2E4ZTQxMjBkYTFlMjcyZGI5L3cvMjAwL2gvMjAwIiwibmFtZSI6Iumds+aYjDAyIiwiZW1haWwiOiJjaGFuZ0BjaGFuZy5jb20iLCJfaWQiOiI1NWMxNzEwZGY5NmJiZTg0NzY4MzI1MmEiLCJpc05ldyI6dHJ1ZSwicmVnaW9uIjoiY24ifX0=; TEAMBITION_SESSIONID.sig=PfTE50ypOxA4uf09mgP9DR2IjKQ")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store.Init(w, r, true)
			session, err := sessions.Get(cookiekey, store)
			assert.Nil(err)
			assert.Equal(float64(1485158874813), session.Values["authUpdated"])
			assert.Equal("55c1710df96bbe847683252a", session.Values["uid"])
			q := session.Values["user"].(map[string]interface{})

			assert.Equal(true, q["isNew"])
			assert.Equal(cookiekey, session.Name())
			assert.NotNil(session.Store())
		})
		handler.ServeHTTP(recorder, req)
	})
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
