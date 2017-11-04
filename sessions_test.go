package sessions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-http-utils/cookie"
	"github.com/go-http-utils/cookie-session"
	"github.com/stretchr/testify/assert"
)

// Session ...
type Session struct {
	*sessions.Meta `json:"-"`
	Name           string `json:"name"`
	Age            int64  `json:"age"`
	Authed         int64  `json:"authed"`
}

// Save ...
func (s *Session) Save() error {
	return s.GetStore().Save(s)
}

func (s *Session) Destroy() error {
	return s.GetStore().Destroy(s)
}

func TestSessions(t *testing.T) {

	SessionName := "teambition"
	NewSessionName := "teambition-new"
	SessionKeys := []string{"keyxxx"}

	t.Run("Sessions use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		store := sessions.New()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroom"
			session.Age = 99
			err = session.Save()
			assert.Nil(err)
			assert.True(session.IsNew())
			assert.True(session.GetSID() == "")
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		store = sessions.New()
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal("mushroom", session.Name)
			assert.Equal(int64(99), session.Age)
			assert.False(session.IsNew())
			assert.True(session.GetSID() != "")
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session=====

		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		store = sessions.New()
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal("mushroom", session.Name)
			assert.Equal(int64(99), session.Age)
			assert.False(session.IsNew())
			assert.True(session.GetSID() != "")
		})
		handler.ServeHTTP(recorder, req)
	})

	t.Run("Sessions with sign session that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		store := sessions.New()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroom"
			session.Age = 99
			session.Save()

			session = &Session{Meta: &sessions.Meta{}}
			store.Load(NewSessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroomnew"
			session.Age = 100
			session.Save()

		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, _ = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		store = sessions.New()
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal("mushroom", session.Name)
			assert.Equal(int64(99), session.Age)

			session = &Session{Meta: &sessions.Meta{}}
			store.Load(NewSessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal("mushroomnew", session.Name)
			assert.Equal(int64(100), session.Age)

		})
		handler.ServeHTTP(recorder, req)

	})

	t.Run("Sessions with Name() and Store()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		store := sessions.New(&sessions.Options{
			Path:     "xxx.com",
			HTTPOnly: true,
			MaxAge:   64,
			Domain:   "ttt.com",
			Secure:   true,
		})

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroom"
			session.Age = 99

			session.Save()

			assert.Equal(SessionName, session.GetName())
			assert.NotNil(session.GetStore())

		})
		handler.ServeHTTP(recorder, req)
		cookies, _ := getCookie(SessionName, recorder)
		assert.Equal("ttt.com", cookies.Domain)
		assert.Equal("xxx.com", cookies.Path)
		assert.Equal(true, cookies.HttpOnly)
		assert.Equal(64, cookies.MaxAge)
		assert.Equal(true, cookies.Secure)

	})

	t.Run("Sessions with Encode() and Decode()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		dt := make(map[string]interface{})
		err := sessions.Decode("xx", dt)
		assert.NotNil(err)

		_, err = sessions.Encode(make(chan int))
		assert.NotNil(err)
	})

	t.Run("Sessions don't override old value when seting same value that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		assert.Nil(err)
		recorder := httptest.NewRecorder()

		store := sessions.New()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroom"
			session.Age = 99

			session.Save()
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		store = sessions.New()
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = "mushroom"
			session.Age = 99
			session.Save()
		})
		handler.ServeHTTP(recorder, req)
	})
}

// User ...
type User struct {
	AvatarURL string `json:"avatarUrl"`
	IsNe      bool   `json:"isNew"`
	ID        string `json:"_id"`
}

// TbSession ...
type TbSession struct {
	*sessions.Meta `json:"-"`
	AuthUpdated    int64  `json:"authUpdated"`
	NextURL        string `json:"nextUrl"`
	TS             int64  `json:"ts"`
	UID            string `json:"uid"`
	User           User   `json:"user"`
}

func TestSessionCompatible(t *testing.T) {
	SessionName := "TEAMBITION_SESSIONID"
	SessionKeys := []string{"tb-accounts"}

	store := sessions.New()

	t.Run("gearsession should be compatible with old session component that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Cookie", "TEAMBITION_SESSIONID=eyJhdXRoVXBkYXRlZCI6MTQ4NTE1ODg3NDgxMywibmV4dFVybCI6Imh0dHA6Ly9wcm9qZWN0LmNpL3Byb2plY3RzIiwidHMiOjE0ODY2MDkzNTA5NjAsInVpZCI6IjU1YzE3MTBkZjk2YmJlODQ3NjgzMjUyYSIsInVzZXIiOnsiYXZhdGFyVXJsIjoiaHR0cDovL3N0cmlrZXIucHJvamVjdC5jaS90aHVtYm5haWwvMDEwa2UyZTMzODQ3ZjQzNzhlY2E4ZTQxMjBkYTFlMjcyZGI5L3cvMjAwL2gvMjAwIiwibmFtZSI6Iumds+aYjDAyIiwiZW1haWwiOiJjaGFuZ0BjaGFuZy5jb20iLCJfaWQiOiI1NWMxNzEwZGY5NmJiZTg0NzY4MzI1MmEiLCJpc05ldyI6dHJ1ZSwicmVnaW9uIjoiY24ifX0=; TEAMBITION_SESSIONID.sig=PfTE50ypOxA4uf09mgP9DR2IjKQ")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &TbSession{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal(int64(1485158874813), session.AuthUpdated)
			assert.Equal("http://project.ci/projects", session.NextURL)
			assert.Equal(int64(1486609350960), session.TS)
			assert.Equal("55c1710df96bbe847683252a", session.UID)

			assert.Equal("http://striker.project.ci/thumbnail/010ke2e33847f4378eca8e4120da1e272db9/w/200/h/200", session.User.AvatarURL)
			assert.Equal("55c1710df96bbe847683252a", session.User.ID)
			assert.Equal(true, session.User.IsNe)
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

func migrateCookies(recorder *httptest.ResponseRecorder, req *http.Request) {
	for _, cookie := range recorder.Result().Cookies() {
		if !cookie.Expires.IsZero() && !cookie.Expires.Before(time.Now()) {
			req.AddCookie(cookie)
		}
	}
}
