package cookiesession

import "net/http"

//Session ...
type Session struct {
	Values map[interface{}]interface{}
	store  Store
	name   string
}

//NewSession ...
func NewSession(name string, store Store) *Session {
	session := &Session{
		Values: make(map[interface{}]interface{}),
		store:  store,
		name:   name,
	}
	return session
}

//Save ....
func (s *Session) Save(r *http.Request, w http.ResponseWriter) {
	s.store.Save(r, w, s)
}

// Store is an interface for custom session stores.
type Store interface {
	// Get should return a cached session.
	Get(r *http.Request, name string) (*Session, error)

	// New should create and return a new session.
	New(r *http.Request, name string) (*Session, error)

	// Save should persist session to the underlying store implementation.
	Save(r *http.Request, w http.ResponseWriter, s *Session) error
}
