package sessions

import (
	"errors"
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

// NewCookieStore an CookieStore instance
func NewCookieStore(Keys []string, options ...*Options) (store *CookieStore) {
	if len(Keys) == 0 {
		panic(errors.New("invalid keys parameter"))
	}
	store = &CookieStore{keys: Keys}
	if len(options) > 0 {
		store.options = options[0]
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
func (c *CookieStore) Init(w http.ResponseWriter, r *http.Request, signed bool) {
	c.cookie = cookie.New(w, r, c.keys)
	c.signed = signed
	return
}

// Get existed session from Request's cookies
func (c *CookieStore) Get(name string) (data map[string]interface{}, err error) {
	val, _ := c.cookie.Get(name, c.signed)
	if val != "" {
		Decode(val, &data)
	}
	return
}

// Save session to Response's cookie
func (c *CookieStore) Save(name string, data map[string]interface{}) (err error) {
	opts := &cookie.Options{
		Path:     "/",
		HTTPOnly: true,
		Signed:   c.signed,
		MaxAge:   24 * 60 * 60,
	}
	if c.options != nil {
		opts.Path = c.options.Path
		opts.Domain = c.options.Domain
		opts.MaxAge = c.options.MaxAge
		opts.Secure = c.options.Secure
		opts.HTTPOnly = c.options.HTTPOnly
	}
	val, err := Encode(data)
	if err == nil {
		c.cookie.Set(name, val, opts)
	}
	return
}
