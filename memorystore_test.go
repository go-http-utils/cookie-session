package sessions_test

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"time"

	"github.com/go-http-utils/cookie"
	"github.com/go-http-utils/cookie-session"
	"github.com/stretchr/testify/assert"
)

var (
	username       = "mushroom"
	useage   int64 = 99

	secondUserName       = "mushroomnew"
	secondUsage    int64 = 100
	store                = sessions.NewMemoryStore()
)

func TestMemoryStore(t *testing.T) {

	SessionName := "teambition"
	NewSessionName := "teambition-new"
	SessionKeys := []string{"keyxxx"}

	t.Run("Sessions use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage
			err = session.Save()
			assert.Nil(err)
			assert.True(session.IsNew())
			assert.True(session.GetSID() == "")
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal(username, session.Name)
			assert.Equal(int64(useage), session.Age)
			assert.False(session.IsNew())
			assert.True(session.GetSID() != "")
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session=====

		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal(username, session.Name)
			assert.Equal(useage, session.Age)
			assert.False(session.IsNew())
			assert.True(session.GetSID() != "")
		})
		handler.ServeHTTP(recorder, req)
	})
	t.Run("Sessions with sign session that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage
			session.Save()

			session = &Session{Meta: &sessions.Meta{}}
			store.Load(NewSessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = secondUserName
			session.Age = secondUsage
			session.Save()

		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, _ = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal(username, session.Name)
			assert.Equal(useage, session.Age)

			session = &Session{Meta: &sessions.Meta{}}
			store.Load(NewSessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal(secondUserName, session.Name)
			assert.Equal(secondUsage, session.Age)

		})
		handler.ServeHTTP(recorder, req)

	})
	t.Run("Sessions with Name() and Store()  that should be", func(t *testing.T) {
		assert := assert.New(t)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		store := sessions.NewMemoryStore(&sessions.Options{
			Path:     "xxx.com",
			HTTPOnly: true,
			MaxAge:   64,
			Domain:   "ttt.com",
			Secure:   true,
		})

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage

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
	t.Run("Sessions donn't override old value when seting same value that should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		assert.Nil(err)
		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage

			session.Save()
		})
		handler.ServeHTTP(recorder, req)

		//====== reuse session =====
		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage
			session.Save()
		})
		handler.ServeHTTP(recorder, req)
	})
	t.Run("Sessions with high goroutine should be", func(t *testing.T) {
		assert := assert.New(t)
		req, err := http.NewRequest("GET", "/", nil)
		assert.Nil(err)
		recorder := httptest.NewRecorder()

		store := sessions.NewMemoryStore(&sessions.Options{
			Path:     "xxx.com",
			HTTPOnly: true,
			MaxAge:   2,
			Domain:   "ttt.com",
			Secure:   true,
		})

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
			session.Name = username
			session.Age = useage
			session.Save()

			var wg sync.WaitGroup
			wg.Add(10000)
			for i := 0; i < 10000; i++ {
				go func() {
					newid := genID()
					sess := &Session{Meta: &sessions.Meta{}}
					store.Load(newid, sess, cookie.New(w, r, SessionKeys...))
					sess.Name = username
					sess.Age = useage
					sess.Save()
					wg.Done()
				}()
			}
			wg.Wait()
		})
		handler.ServeHTTP(recorder, req)
		time.Sleep(time.Second * 3)
		assert.Equal(0, store.Len())
		//====== reuse session =====
		req, err = http.NewRequest("GET", "/", nil)
		migrateCookies(recorder, req)

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := &Session{Meta: &sessions.Meta{}}
			store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))

			assert.Equal("", session.Name)
			assert.Equal(int64(0), session.Age)
		})
		handler.ServeHTTP(recorder, req)

		store.Destroy()
	})
}
func genID() string {
	buf := make([]byte, 12)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}
