package sessions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// Session stores the values and optional configuration for a session.
type Session struct {
	Values map[string]interface{}
	SID    string
	store  Store
	name   string
	cache  map[string]interface{}
	w      http.ResponseWriter
	req    *http.Request
}

// NewSession to create new session instance
func NewSession(name string, store Store, w http.ResponseWriter, r *http.Request) (session *Session) {
	session = &Session{
		Values: make(map[string]interface{}),
		store:  store,
		name:   name,
		w:      w,
		req:    r,
		cache:  make(map[string]interface{}),
	}
	return
}

// Name returns the name used to register the session
func (s *Session) Name() string {
	return s.name
}

// Store returns the session store used to register the session
func (s *Session) Store() Store {
	return s.store
}

// AddCache to add data cache for thirdparty store implement, like store the last value to check whether changed when saving
func (s *Session) AddCache(key string, val interface{}) {
	switch v := val.(type) {
	case map[string]interface{}:
		s.cache[key] = copyMap(val.(map[string]interface{}))
	default:
		s.cache[key] = v
	}
}

func copyMap(cache map[string]interface{}) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for k, v := range cache {
		switch val := v.(type) {
		case map[string]interface{}:
			data[k] = copyMap(val)
		default:
			data[k] = val
		}
	}
	return
}

// GetCache cache data from current Session
func (s *Session) GetCache(name string) interface{} {
	return s.cache[name]
}

// Save is a convenience method to save current session
func (s *Session) Save() (err error) {
	s.store.Save(s.w, s.req, s)
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
	Get(name string, w http.ResponseWriter, r *http.Request) (session *Session, err error)

	// Save should persist session to the underlying store implementation.
	Save(w http.ResponseWriter, r *http.Request, session *Session) error
}
