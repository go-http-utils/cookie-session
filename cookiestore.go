package cookiesession

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/http"

	"github.com/go-http-utils/cookie"
)

//New an CookieStore instance
func New(opts ...*cookie.GlobalOptions) Store {
	var opt *cookie.GlobalOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	cs := &CookieStore{
		options: opt,
	}
	return cs
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	options *cookie.GlobalOptions
}

//Get existed session from Request's cookies
func (c *CookieStore) Get(r *http.Request, name string) (session *Session, err error) {
	session, err = c.New(r, name)
	session.name = name
	session.store = c
	return
}

//New an session instance
func (c *CookieStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(name, c)
	var err error
	cookies := cookie.New(nil, r, c.options)
	val, err := cookies.Get(name)
	if val != "" {
		c.Decode(val, &session.Values)
	}
	return session, err
}

//Save session to Response's cookie
func (c *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	encoded, err := c.Encode(s.Values)
	if err != nil {
		return err
	}
	cookies := cookie.New(w, r, c.options)
	cookies.Set(s.name, encoded)
	return nil
}

//Encode the value by Serializer and Base64
func (c *CookieStore) Encode(value interface{}) (string, error) {

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

//Decode the value to dst .
func (c *CookieStore) Decode(value string, dst interface{}) error {

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
