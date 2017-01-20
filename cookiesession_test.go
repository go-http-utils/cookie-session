package cookiesession

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGearsession(t *testing.T) {

	t.Run("gearsession use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/health-check", nil)
		store := New()

		cookiekey := "teambition"
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session, _ := store.Get(r, cookiekey)
			session.Values["name"] = "mushroom"
			session.Values[66] = 99
			session.Save(r, w)

		})
		handler.ServeHTTP(recorder, req)

		cookies, err := getCookie(cookiekey, recorder)
		assert.Nil(err)
		assert.NotNil(cookies.Value)
		t.Log(cookies.Value)
		//======reuse=====
		req, err = http.NewRequest("GET", "/health-check", nil)
		store = New(nil)
		req.AddCookie(cookies)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, cookiekey)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(99, session.Values[66])
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		req, err = http.NewRequest("GET", "/health-check", nil)
		req.AddCookie(cookies)
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, cookiekey)
			assert.Equal("mushroom", session.Values["name"])
			assert.Equal(99, session.Values[66])
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
