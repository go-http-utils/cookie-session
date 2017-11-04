package sessions

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-http-utils/cookie"
)

// Version is this package's version
const Version = "1.0.0"

// Store is an interface for custom session stores.
type Store interface {
	// Load should load data from cookie and store, set it into session instance.
	// error indicates that session validation failed, or other thing.
	// Sessions.Init should be called in Load, even if error occured.
	Load(name string, session Sessions, cookie *cookie.Cookies) error
	// Save should persist session to the underlying store implementation.
	Save(session Sessions) error
	// Destroy should destroy the session.
	Destroy(session Sessions) error
}

// Sessions ...
type Sessions interface {
	// Init sets current cookie.Cookies and Store to the session instance.
	Init(name, sid string, c *cookie.Cookies, store Store, lastValue string)
	// GetSID returns the session' sid
	GetSID() string
	// GetName returns the session' name
	GetName() string
	// GetStore returns the session' store
	GetStore() Store
	// GetCookie returns the session' cookie
	GetCookie() *cookie.Cookies
	// IsChanged to check current session's value whether is changed
	IsChanged(val string) bool
	// IsNew to check the current session whether it's new user
	IsNew() bool
}

// Meta stores the values and optional configuration for a session.
type Meta struct {
	// Values map[string]interface{}
	sid       string
	store     Store
	name      string
	cookie    *cookie.Cookies
	lastValue string
}

// Init sets current cookie.Cookies and Store to the session instance.
func (s *Meta) Init(name, sid string, c *cookie.Cookies, store Store, lastValue string) {
	s.name = name
	s.sid = sid
	s.cookie = c
	s.store = store
	s.lastValue = lastValue
}

// GetSID returns the session' sid
func (s *Meta) GetSID() string {
	return s.sid
}

// GetName returns the name used to register the session
func (s *Meta) GetName() string {
	return s.name
}

// GetStore returns the session store used to register the session
func (s *Meta) GetStore() Store {
	return s.store
}

// GetCookie returns the session' cookie
func (s *Meta) GetCookie() *cookie.Cookies {
	return s.cookie
}

// IsChanged to check current session's value whether is changed
func (s *Meta) IsChanged(val string) bool {
	return s.lastValue != val
}

// IsNew to check the current session whether it's new user
func (s *Meta) IsNew() bool {
	return s.sid == ""
}

// Encode the value by Serializer and Base64
func Encode(value interface{}) (str string, err error) {
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
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
