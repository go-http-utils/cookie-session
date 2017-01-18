package cookiesession

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/http"

	"github.com/go-http-utils/cookie"
)

//NewCookieStore ...
func NewCookieStore(opts *cookie.Options) Store {
	cs := &CookieStore{
		options: opts,
	}
	return cs
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	options *cookie.Options
}

//Get ...
func (c *CookieStore) Get(r *http.Request, name string) (session *Session, err error) {
	session, err = c.New(r, name)
	session.name = name
	session.store = c
	return
}

//New ...
func (c *CookieStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(name, c)
	var err error
	cookies := cookie.New(nil, r, c.options)
	val, err := cookies.Get(name, nil)
	if val != "" {
		c.Decode(name, val, &session.Values)
	}
	return session, err
}

//Save ...
func (c *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	encoded, err := c.Encode(s.name, s.Values)
	if err != nil {
		return err
	}
	cookies := cookie.New(w, r, c.options)
	cookies.Set(s.name, encoded, nil)
	return nil
}

//Encode ...
func (c *CookieStore) Encode(name string, value interface{}) (string, error) {

	//Serializer
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	b := buf.Bytes()

	//base64
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(b)))
	base64.URLEncoding.Encode(encoded, b)

	return string(encoded), nil
}

//Decode decodes a cookie value.
func (c *CookieStore) Decode(name, value string, dst interface{}) error {

	//base64
	b, err := decode([]byte(value))
	if err != nil {
		return err
	}

	//Serializer
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

// decode decodes a cookie using base64.
func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}
