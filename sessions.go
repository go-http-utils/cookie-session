package sessions

import (
	"encoding/base64"
	"encoding/json"
)

// Session stores the values and optional configuration for a session.
type Session struct {
	SID    string
	Values map[string]interface{}
	store  Store
	name   string
	oldval string
}

// Name returns the name used to register the session
func (s *Session) Name() string {
	return s.name
}

// Store returns the session store used to register the session
func (s *Session) Store() Store {
	return s.store
}

// Get a session instance by name and any kind of stores
func Get(name string, store Store) (session *Session, err error) {
	session = &Session{
		Values: make(map[string]interface{}),
		store:  store,
		name:   name,
	}
	val, err := store.Get(name)
	if val != "" {
		session.oldval = val
		err = Decode(val, &session.Values)
	}
	return
}

// New a new session instance by name
func (s *Session) New(name string) *Session {
	session := &Session{
		Values: make(map[string]interface{}),
		store:  s.store,
		name:   name,
	}
	return session
}

// Get a new session instance by name base on current store
func (s *Session) Get(name string) *Session {
	val, _ := s.store.Get(name)
	if val != "" {
		s.oldval = val
		Decode(val, &s.Values)
	}
	return s
}

// Save is a convenience method to save current session
func (s *Session) Save() (err error) {
	encoded, err := Encode(s.Values)
	if err == nil && s.oldval != encoded {
		s.store.Save(s.name, encoded)
	}
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
	// Get should return a cached session string.
	Get(name string) (string, error)

	// Save should persist session to the underlying store implementation.
	Save(name string, data string) error
}
