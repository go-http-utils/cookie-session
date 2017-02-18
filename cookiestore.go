package sessions

import "github.com/go-http-utils/cookie"

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	opts *cookie.Options
}

// New returns a CookieStore store
func New(opts *cookie.Options) *CookieStore {
	// always use signed cookie for CookieStore session
	opts.Signed = true
	return &CookieStore{opts}
}

// Load loads session values from request cookie
// session will always be initialized with Sessions Init method.
func (s *CookieStore) Load(name string, session Sessions, c *cookie.Cookies) error {
	val, err := c.Get(name, s.opts.Signed)
	if val != "" {
		err = Decode(val, session)
	}
	session.Init(name, val, c, s)
	return err
}

// Save saves session to Response's cookie
func (s *CookieStore) Save(session Sessions) error {
	val, err := Encode(session)
	if err == nil {
		session.GetCookie().Set(session.GetName(), val, s.opts)
	}
	return err
}
