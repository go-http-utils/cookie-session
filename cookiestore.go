package sessions

import "github.com/go-http-utils/cookie"

// Options stores configuration for a session or session store.
//
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// New returns an CookieStore instance
func New(options ...*Options) (store *CookieStore) {
	opts := &cookie.Options{
		Path:     "/",
		HTTPOnly: true,
		Signed:   true,
		MaxAge:   24 * 60 * 60,
	}
	if len(options) > 0 && options[0] != nil {
		temp := options[0]
		opts.Path = temp.Path
		opts.Domain = temp.Domain
		opts.MaxAge = temp.MaxAge
		opts.Secure = temp.Secure
		opts.HTTPOnly = temp.HTTPOnly
	}
	store = &CookieStore{opts}
	return
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	opts *cookie.Options
}

// Load a session by name and any kind of stores
func (c *CookieStore) Load(name string, session Sessions, cookie *cookie.Cookies) error {
	val, err := cookie.Get(name, c.opts.Signed)
	if val != "" {
		err = Decode(val, &session)
	}
	// should call Init even if err
	session.Init(name, val, cookie, c, val)
	return err
}

// Save session to Response's cookie
func (c *CookieStore) Save(session Sessions) (err error) {
	val, err := Encode(session)
	if err == nil && session.IsChanged(val) {
		session.GetCookie().Set(session.GetName(), val, c.opts)
	}
	return
}

// Destroy destroy the session
func (c *CookieStore) Destroy(session Sessions) (err error) {
	session.GetCookie().Remove(session.GetName(), c.opts)
	return
}
