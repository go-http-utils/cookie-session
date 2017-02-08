package cookiesession

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"time"
)

// Options stores configuration for a session or session store.
//
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	Signed   bool // optional
}

// Session stores the values and optional configuration for a session.
type Session struct {
	SID     string
	Values  map[interface{}]interface{}
	Options *Options
	store   Store
	name    string
}

// Name returns the name used to register the session
func (s *Session) Name() string {
	return s.name
}

// Store returns the session store used to register the session
func (s *Session) Store() Store {
	return s.store
}

// NewSession is called by session stores to create a new session instance
func NewSession(name string, store Store) *Session {
	session := &Session{
		Values: make(map[interface{}]interface{}),
		store:  store,
		name:   name,
	}
	return session
}

// Save is a convenience method to save current session
func (s *Session) Save(options ...*Options) {
	s.store.Save(s, options...)
}

// Encode the value by Serializer and Base64
func Encode(value interface{}) (string, error) {

	//Serializer
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	b := buf.Bytes()

	//Base64
	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(encoded, b)

	return string(encoded), nil
}

// Decode the value to dst .
func Decode(value string, dst interface{}) error {

	//base64
	val := []byte(value)
	decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(val)))
	b, err := base64.RawURLEncoding.Decode(decoded, val)
	if err != nil {
		return err
	}
	//Serializer
	dec := gob.NewDecoder(bytes.NewBuffer(decoded[:b]))
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

// NewCookie returns an http.Cookie with the options set. It also sets
// the Expires field calculated based on the MaxAge value, for Internet
// Explorer compatibility.
func NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HTTPOnly,
	}
	if options.MaxAge > 0 {
		d := time.Duration(options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if options.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}

// Store is an interface for custom session stores.
type Store interface {
	// Get should return a cached session.
	Get(name string, signed ...bool) (*Session, error)

	// New should create and return a new session.
	New(name string) (*Session, error)

	// Save should persist session to the underlying store implementation.
	Save(s *Session, options ...*Options) error
}
