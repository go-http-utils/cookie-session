package sessions

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-http-utils/cookie"
	"github.com/stretchr/testify/assert"
)

type Session struct {
	*Meta  `json:"-"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Authed int64  `json:"authed"`
}

func (s *Session) IsNew() bool {
	return s.GetSID() == ""
}

func (s *Session) Save() error {
	return s.GetStore().Save(s)
}

func TestSessionWithCookieStore(t *testing.T) {
	SessionName := "Sess"
	SessionKeys := []string{"keyxxx"}
	store := New(&cookie.Options{})

	t.Run("Sessions use default options that should be", func(t *testing.T) {
		assert := assert.New(t)

		req, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys))
			assert.Equal("", session.UserID)
			assert.Equal("", session.Name)
			assert.Equal(int64(0), session.Authed)
			assert.True(session.IsNew())

			session.UserID = "user123"
			session.Name = "test"
			session.Authed = time.Now().Unix() + 100
			assert.Nil(session.Save())
		})
		handler.ServeHTTP(recorder, req)

		//====== load session =====
		req, _ = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &Meta{}}

			err := store.Load(SessionName, session, cookie.New(w, r, SessionKeys))
			assert.Nil(err)
			assert.Equal("user123", session.UserID)
			assert.Equal("test", session.Name)
			assert.True(session.Authed > time.Now().Unix())
			assert.False(session.IsNew())
		})
		handler.ServeHTTP(recorder, req)

		//====== load session with wrong key =====
		req, _ = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &Meta{}}

			err := store.Load(SessionName, session, cookie.New(w, r, []string{"keyxxx1"}))
			assert.NotNil(err)
			assert.Equal("", session.UserID)
			assert.Equal("", session.Name)
			assert.True(session.IsNew())
		})
		handler.ServeHTTP(recorder, req)

		//======reuse=====
		// store = sessions.NewCookieStore([]string{})
		// req, err = http.NewRequest("GET", "/", nil)
		// req.AddCookie(cookies)
		// handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 	session, _ := store.Get(cookiekey, w, r)

		// 	assert.Equal("mushroom", session.Values["name"])
		// 	assert.Equal(float64(99), session.Values["num"])
		// })
		// handler.ServeHTTP(recorder, req)
	})

	// t.Run("Sessions with sign session that should be", func(t *testing.T) {
	// 	assert := assert.New(t)
	// 	recorder := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("GET", "/", nil)

	// 	store := sessions.NewCookieStore([]string{"key"})

	// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		session, err := store.Get(cookiekey, w, r)
	// 		session.Values["name"] = "mushroom"
	// 		session.Values["num"] = 99
	// 		session.Save()
	// 		assert.Nil(err)
	// 		session, err = store.Get(cookieNewKey, w, r)
	// 		session.Values["name"] = "teambition-n"
	// 		session.Values["num"] = 100
	// 		session.Save()
	// 		assert.Nil(err)
	// 	})
	// 	handler.ServeHTTP(recorder, req)

	// 	//======reuse=====
	// 	store = sessions.NewCookieStore([]string{"key"})
	// 	req, _ = http.NewRequest("GET", "/", nil)
	// 	cookies, _ := getCookie(cookiekey, recorder)
	// 	req.AddCookie(cookies)
	// 	cookies, _ = getCookie(cookiekey+".sig", recorder)
	// 	req.AddCookie(cookies)

	// 	cookies, _ = getCookie(cookieNewKey, recorder)
	// 	req.AddCookie(cookies)
	// 	cookies, _ = getCookie(cookieNewKey+".sig", recorder)
	// 	req.AddCookie(cookies)

	// 	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		session, err := store.Get(cookiekey, w, r)

	// 		assert.Nil(err)
	// 		assert.Equal("mushroom", session.Values["name"])
	// 		assert.Equal(float64(99), session.Values["num"])

	// 		session, err = store.Get(cookieNewKey, w, r)
	// 		assert.Nil(err)
	// 		assert.Equal("teambition-n", session.Values["name"])
	// 		assert.Equal(float64(100), session.Values["num"])

	// 		session, err = store.Get(cookieNewKey+"new", w, r)
	// 		assert.Nil(err)
	// 		assert.Equal(0, len(session.Values))

	// 	})
	// 	handler.ServeHTTP(recorder, req)

	// })
	// t.Run("Sessions with Name() and Store()  that should be", func(t *testing.T) {
	// 	assert := assert.New(t)
	// 	recorder := httptest.NewRecorder()
	// 	req, _ := http.NewRequest("GET", "/", nil)

	// 	store := sessions.NewCookieStore([]string{"key"}, &sessions.Options{
	// 		Path:     "xxx.com",
	// 		HTTPOnly: true,
	// 		MaxAge:   64,
	// 		Domain:   "ttt.com",
	// 		Secure:   true,
	// 	})

	// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// 		session, err := store.Get(cookiekey, w, r)
	// 		session.Values["name"] = "mushroom"
	// 		session.Values["num"] = 99
	// 		session.Save()

	// 		assert.Nil(err)
	// 		assert.Equal(cookiekey, session.Name())
	// 		assert.NotNil(session.Store())

	// 	})
	// 	handler.ServeHTTP(recorder, req)
	// 	cookies, _ := getCookie(cookiekey, recorder)
	// 	assert.Equal("ttt.com", cookies.Domain)
	// 	assert.Equal("xxx.com", cookies.Path)
	// 	assert.Equal(true, cookies.HttpOnly)
	// 	assert.Equal(64, cookies.MaxAge)
	// 	assert.Equal(true, cookies.Secure)

	// })

	// t.Run("Sessions with Encode() and Decode()  that should be", func(t *testing.T) {
	// 	assert := assert.New(t)
	// 	dt := make(map[string]interface{})
	// 	err := sessions.Decode("xx", dt)
	// 	assert.NotNil(err)

	// 	_, err = sessions.Encode(make(chan int))
	// 	assert.NotNil(err)
	// })

	// t.Run("Sessions donn't override old value when seting same value that should be", func(t *testing.T) {
	// 	assert := assert.New(t)
	// 	req, err := http.NewRequest("GET", "/", nil)
	// 	recorder := httptest.NewRecorder()

	// 	store := sessions.NewCookieStore([]string{})

	// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		session, err := store.Get(cookiekey, w, r)
	// 		session.Values["name"] = "mushroom"
	// 		session.Values["num"] = 99
	// 		session.Save()
	// 		assert.Nil(err)
	// 	})
	// 	handler.ServeHTTP(recorder, req)

	// 	//======reuse=====
	// 	store = sessions.NewCookieStore([]string{})
	// 	cookies, err := getCookie(cookiekey, recorder)
	// 	assert.Nil(err)
	// 	assert.NotNil(cookies.Value)
	// 	req, err = http.NewRequest("GET", "/", nil)

	// 	req.AddCookie(cookies)

	// 	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		session, err := store.Get(cookiekey, w, r)
	// 		session.Save()
	// 		assert.Nil(err)
	// 	})
	// 	handler.ServeHTTP(recorder, req)
	// })

}

// func TestSessionCompatible(t *testing.T) {
// 	cookiekey := "TEAMBITION_SESSIONID"
// 	store := sessions.NewCookieStore([]string{"tb-accounts"})

// 	t.Run("gearsession should be compatible with old session component that should be", func(t *testing.T) {
// 		assert := assert.New(t)
// 		recorder := httptest.NewRecorder()

// 		req, _ := http.NewRequest("GET", "/", nil)
// 		req.Header.Set("Cookie", "TEAMBITION_SESSIONID=eyJhdXRoVXBkYXRlZCI6MTQ4NTE1ODg3NDgxMywibmV4dFVybCI6Imh0dHA6Ly9wcm9qZWN0LmNpL3Byb2plY3RzIiwidHMiOjE0ODY2MDkzNTA5NjAsInVpZCI6IjU1YzE3MTBkZjk2YmJlODQ3NjgzMjUyYSIsInVzZXIiOnsiYXZhdGFyVXJsIjoiaHR0cDovL3N0cmlrZXIucHJvamVjdC5jaS90aHVtYm5haWwvMDEwa2UyZTMzODQ3ZjQzNzhlY2E4ZTQxMjBkYTFlMjcyZGI5L3cvMjAwL2gvMjAwIiwibmFtZSI6Iumds+aYjDAyIiwiZW1haWwiOiJjaGFuZ0BjaGFuZy5jb20iLCJfaWQiOiI1NWMxNzEwZGY5NmJiZTg0NzY4MzI1MmEiLCJpc05ldyI6dHJ1ZSwicmVnaW9uIjoiY24ifX0=; TEAMBITION_SESSIONID.sig=PfTE50ypOxA4uf09mgP9DR2IjKQ")
// 		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			session, err := store.Get(cookiekey, w, r)
// 			assert.Nil(err)
// 			assert.Equal(float64(1485158874813), session.Values["authUpdated"])
// 			assert.Equal("55c1710df96bbe847683252a", session.Values["uid"])
// 			q := session.Values["user"].(map[string]interface{})

// 			assert.Equal(true, q["isNew"])
// 			assert.Equal(cookiekey, session.Name())
// 			assert.NotNil(session.Store())
// 		})
// 		handler.ServeHTTP(recorder, req)
// 	})
// }
// func getCookie(name string, recorder *httptest.ResponseRecorder) (*http.Cookie, error) {
// 	var err error
// 	res := &http.Response{Header: http.Header{"Set-Cookie": recorder.HeaderMap["Set-Cookie"]}}
// 	for _, val := range res.Cookies() {
// 		if val.Name == name {
// 			return val, nil
// 		}
// 	}
// 	return nil, err
// }

func migrateCookies(recorder *httptest.ResponseRecorder, req *http.Request) {
	for _, cookie := range recorder.Result().Cookies() {
		req.AddCookie(cookie)
	}
}

// type MyStore map[string]string

// func (s MyStore) Load(session Sessions, c *cookie.Cookies) error {
// 	name := session.GetName()
// 	sid, err := c.Get(name)
// 	if sid != "" {
// 		if val := s[name+":"+sid]; val != "" {
// 			err = Decode(val, session)
// 		} else {
// 			sid = ""
// 		}
// 	}

// 	session.Init(c, s, sid)
// 	if sid == "" && err == nil {
// 		err = errors.New("session not exists")
// 	}
// 	return err
// }

// func (s MyStore) Save(session Sessions) error {
// 	val, err := Encode(session)
// 	fmt.Println(1111, val, err)
// 	if err == nil {
// 		s[session.GetName()] = val
// 	}
// 	return err
// }
