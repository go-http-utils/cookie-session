package sessions

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/go-http-utils/cookie"
)

// Store is an interface for custom session stores.
type Store interface {
	// Save should persist session to the underlying store implementation.
	Save(session Sessions) error
}

// Sessions ...
type Sessions interface {
	// Init sets current cookie.Cookies and Store to the session instance.
	Init(name, sid string, c *cookie.Cookies, store Store)
	// GetSID returns the session' sid
	GetSID() string
	// GetName returns the session' name
	GetName() string
	// GetStore returns the session' store
	GetStore() Store
	// GetCookie returns the session' cookie
	GetCookie() *cookie.Cookies
}

// Meta implements Sessions interface.
// You can define a custom Session type and embed Meta type:
//
//  type Session struct {
//  	*Meta  `json:"-"`
//  	UserID string `json:"userId"`
//  	Name   string `json:"name"`
//  	Authed int64  `json:"authed"`
//  }
//
//  func (s *Session) IsNew() bool {
//  	return s.GetSID() == ""
//  }
//
//  func (s *Session) Save() error {
//  	return s.GetStore().Save(s)
//  }
//
type Meta struct {
	name, sid string
	store     Store
	cookie    *cookie.Cookies
}

// GetName implements Session interface
func (m *Meta) GetName() string {
	return m.name
}

// GetSID implements Session interface
func (m *Meta) GetSID() string {
	return m.sid
}

// GetStore implements Session interface
func (m *Meta) GetStore() Store {
	return m.store
}

// GetCookie implements Session interface
func (m *Meta) GetCookie() *cookie.Cookies {
	return m.cookie
}

// Init implements Session interface
func (m *Meta) Init(name, sid string, c *cookie.Cookies, store Store) {
	m.name = name
	m.sid = sid
	m.cookie = c
	m.store = store
}

// Encode the session to string by Serializer and Base64
func Encode(session interface{}) (string, error) {
	b, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode the string value to session.
func Decode(val string, session interface{}) error {
	b, err := base64.StdEncoding.DecodeString(val)
	if err == nil {
		err = json.Unmarshal(b, session)
	}
	return err
}

// RandSID ...
func RandSID(seed string) string {
	hasher := sha1.New()
	hasher.Write([]byte(seed))

	buf := make([]byte, 16)
	io.ReadFull(rand.Reader, buf)
	hasher.Write(buf)
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
}
