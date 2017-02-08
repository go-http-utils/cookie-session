package cookiesession

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGearsession(t *testing.T) {
	cookiekey := "teambition"
	cookieNewKey := "teambition-new"
	t.Run("gearsession use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/health-check", nil)

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r)
			session, _ := store.Get(cookiekey)
			session.Values["name"] = "mushroom"
			session.Values[66] = 99
			session.Save()

		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		cookies, err := getCookie(cookiekey, recorder)
		assert.Nil(err)
		assert.NotNil(cookies.Value)
		t.Log(cookies.Value)
		req, err = http.NewRequest("GET", "/health-check", nil)

		req.AddCookie(cookies)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r)
			session, _ := store.Get(cookiekey)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(99, session.Values[66])
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		req, err = http.NewRequest("GET", "/health-check", nil)
		req.AddCookie(cookies)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r)
			session, _ := store.Get(cookiekey)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(99, session.Values[66])
		})
		handler.ServeHTTP(recorder, req)
	})
	t.Run("gearsession with New session that should be", func(t *testing.T) {
		assert := assert.New(t)

		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/health-check", nil)
		securekey := []string{"key"}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r, securekey)
			session, _ := store.Get(cookiekey)
			session.Values["name"] = "mushroom"
			session.Values[66] = 99
			session.Save()

			session, _ = store.New(cookieNewKey)
			session.Values["name"] = "teambition-n"
			session.Values[66] = 100
			session.Save(&Options{Path: "/", HTTPOnly: true, Signed: true})
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====

		req, _ = http.NewRequest("GET", "/health-check", nil)

		cookies, _ := getCookie(cookiekey, recorder)
		req.AddCookie(cookies)
		cookies, _ = getCookie(cookieNewKey, recorder)
		req.AddCookie(cookies)
		cookies, _ = getCookie(cookieNewKey+".sig", recorder)
		req.AddCookie(cookies)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r, securekey)
			session, _ := store.Get(cookiekey)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(99, session.Values[66])

			session, _ = store.Get(cookieNewKey, true)
			assert.Equal("teambition-n", session.Values["name"])
			assert.Equal(100, session.Values[66])

		})
		handler.ServeHTTP(recorder, req)

	})
	t.Run("gearsession with Name() and Store()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/health-check", nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := New(w, r)
			session, _ := store.Get(cookiekey)
			session.Values["name"] = "mushroom"
			session.Values[66] = 99
			session.Save()

			assert.Equal(cookiekey, session.Name())
			assert.NotNil(session.Store())
		})
		handler.ServeHTTP(recorder, req)

	})
	t.Run("gearsession with Name() and Store()  that should be", func(t *testing.T) {
		assert := assert.New(t)

		cookies := NewCookie("key", "val", &Options{
			MaxAge:   1,
			Domain:   "tb.com",
			Path:     "/",
			Secure:   true,
			HTTPOnly: true,
		})
		assert.Equal(cookies.Name, "key")
		assert.Equal(cookies.Value, "val")
		assert.Equal(cookies.MaxAge, 1)
		assert.Equal(cookies.Domain, "tb.com")
		assert.Equal(cookies.Path, "/")
		assert.Equal(cookies.HttpOnly, true)
		assert.Equal(cookies.Secure, true)
		assert.NotNil(cookies.Expires)

		cookies = NewCookie("key", "val", &Options{
			MaxAge:   -1,
			Domain:   "tb.com",
			Path:     "/",
			Secure:   true,
			HTTPOnly: true,
		})
		assert.Equal(cookies.MaxAge, -1)
		assert.NotNil(cookies.Expires)
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
