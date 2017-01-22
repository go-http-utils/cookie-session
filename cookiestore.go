package cookiesession

import (
	"net/http"

	"github.com/go-http-utils/cookie"
)

//New an CookieStore instance
func New(opts ...*cookie.GlobalOptions) Store {
	var opt *cookie.GlobalOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &cookie.GlobalOptions{
			MaxAge:   86400 * 7,
			Secure:   true,
			HTTPOnly: true,
			Path:     "/",
		}
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
	session.Options = &Options{
		Path:     c.options.Path,
		Domain:   c.options.Domain,
		MaxAge:   c.options.MaxAge,
		Secure:   c.options.Secure,
		HTTPOnly: c.options.HTTPOnly,
	}
	var err error
	cookies := cookie.New(nil, r, c.options)
	val, err := cookies.Get(name)
	if val != "" {
		Decode(val, &session.Values)
	}
	return session, err
}

//Save session to Response's cookie
func (c *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {
	encoded, err := Encode(s.Values)
	if err != nil {
		return err
	}
	cookies := cookie.New(w, r, c.options)
	cookies.Set(s.name, encoded)
	return nil
}
