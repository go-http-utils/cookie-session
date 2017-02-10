package sessions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// Session stores the values and optional configuration for a session.
type Session struct {
	sid    string
	Values map[string]interface{}
	store  Store
	name   string
	//oldval string
}

// Name returns the name used to register the session
func (s *Session) Name() string {
	return s.name
}

// Store returns the session store used to register the session
func (s *Session) Store() Store {
	return s.store
}

// New a session instance by name and any kind of stores
func New(name string, store Store, w http.ResponseWriter, r *http.Request, signed ...bool) (session *Session, err error) {
	session = &Session{
		store: store,
		name:  name,
	}
	var issigned = false
	if len(signed) > 0 {
		issigned = signed[0]
	}
	store.Init(w, r, issigned)
	session.Values, err = store.Get(name)

	if session.Values == nil {
		session.Values = make(map[string]interface{})
	}
	return
}

// New a new session instance by name base on current store
func (s *Session) New(name string) *Session {
	session := &Session{
		Values: make(map[string]interface{}),
		store:  s.store,
		name:   name,
	}
	return session
}

// Get a new session instance by name base on current store
func (s *Session) Get(name string) (session *Session, err error) {
	session = &Session{

		Values: make(map[string]interface{}),
		store:  s.store,
		name:   name,
	}
	session.Values, err = s.store.Get(session.name)
	return
}

// Save is a convenience method to save current session
func (s *Session) Save() (err error) {
	s.store.Save(s.name, s.Values)
	return
}

// Encode the value by Serializer and Base64
func Encode(value interface{}) (str string, err error) {

	//Serializer
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
	//Base64
	str = base64.StdEncoding.EncodeToString(b)
	return
}

// Decode the value to dst .
func Decode(value string, dst interface{}) (err error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, dst)
	return
}

// Store is an interface for custom session stores.
type Store interface {
	Init(w http.ResponseWriter, r *http.Request, signed bool)

	// Get should return a cached session string.
	Get(name string) (map[string]interface{}, error)

	// Save should persist session to the underlying store implementation.
	Save(name string, data map[string]interface{}) error
}
