package sessions

import (
	"net/http"

	"github.com/go-http-utils/cookie"
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
}

// New an CookieStore instance
func New(Keys ...[]string) (store *CookieStore) {
	if len(Keys) > 0 && len(Keys[0]) > 0 {
		store = &CookieStore{keys: Keys[0]}
	} else {
		store = &CookieStore{}
	}
	return
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	cookie  *cookie.Cookies
	options *Options
	keys    []string
	signed  bool // optional
}

// Init an CookieStore instance
func (c *CookieStore) Init(w http.ResponseWriter, r *http.Request, signed bool, options ...*Options) {
	if len(options) > 0 {
		c.options = options[0]
	}
	if len(c.keys) > 0 && len(c.keys[0]) > 0 {
		c.cookie = cookie.New(w, r, c.keys)
	} else {
		c.cookie = cookie.New(w, r)
	}
	c.signed = signed
	return
}

// Get existed session from Request's cookies
func (c *CookieStore) Get(name string) (val string, err error) {
	val, err = c.cookie.Get(name, c.signed)
	return
}

// Save session to Response's cookie
func (c *CookieStore) Save(name string, data string) error {
	opts := &cookie.Options{
		Path:     "/",
		HTTPOnly: true,
		Signed:   c.signed,
		MaxAge:   24 * 60 * 60,
	}
	opts.Signed = c.signed
	if c.options != nil {
		opts.Path = c.options.Path
		opts.Domain = c.options.Domain
		opts.MaxAge = c.options.MaxAge
		opts.Secure = c.options.Secure
		opts.HTTPOnly = c.options.HTTPOnly
	}
	c.cookie.Set(name, data, opts)
	return nil
}
