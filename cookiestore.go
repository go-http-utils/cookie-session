package cookiesession

import (
	"net/http"

	"github.com/go-http-utils/cookie"
)

// New an CookieStore instance
func New(w http.ResponseWriter, r *http.Request, keys ...[]string) Store {
	cs := &CookieStore{cookie: cookie.New(w, r, keys...)}
	return cs
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	cookie *cookie.Cookies
}

// Get existed session from Request's cookies
func (c *CookieStore) Get(name string, signed ...bool) (session *Session, err error) {
	session = NewSession(name, c)

	val, err := c.cookie.Get(name, signed...)
	if val != "" {
		Decode(val, &session.Values)
	}
	session.name = name
	session.store = c
	return
}

// New an session instance
func (c *CookieStore) New(name string) (session *Session, err error) {
	session = NewSession(name, c)
	session.name = name
	session.store = c
	return
}

// Save session to Response's cookie
func (c *CookieStore) Save(s *Session, options ...*Options) error {
	encoded, err := Encode(s.Values)
	if err != nil {
		return err
	}
	if len(options) > 0 {
		option := options[0]
		opts := &cookie.Options{
			Path:     option.Path,
			Domain:   option.Domain,
			MaxAge:   option.MaxAge,
			Secure:   option.Secure,
			HTTPOnly: option.HTTPOnly,
			Signed:   option.Signed,
		}
		c.cookie.Set(s.name, encoded, opts)
	} else {
		c.cookie.Set(s.name, encoded)
	}
	return nil
}
