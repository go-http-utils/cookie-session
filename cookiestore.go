package sessions

import (
	"net/http"
	"reflect"

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
func NewCookieStore(keys []string, options ...*Options) (store *CookieStore) {
	store = &CookieStore{keys: keys}
	if len(options) > 0 {
		store.options = options[0]
	}
	if len(keys) > 0 && len(keys[0]) > 0 {
		store.signed = true
	}
	return
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	options *Options
	keys    []string
	signed  bool
}

// Get a session instance by name and any kind of stores
func (c *CookieStore) Get(name string, w http.ResponseWriter, r *http.Request) (session *Session, err error) {
	cookie := cookie.New(w, r, c.keys)
	session = NewSession(name, c, w, r)
	val, _ := cookie.Get(name, c.signed)
	if val != "" {
		Decode(val, &session.Values)
	}
	session.AddCache("cookie", cookie)
	session.AddCache("lastvalue", session.Values)
	return
}

// Save session to Response's cookie
func (c *CookieStore) Save(w http.ResponseWriter, r *http.Request, session *Session) (err error) {
	if reflect.DeepEqual(session.GetCache("lastvalue"), session.Values) {
		return
	}
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
	val, err := Encode(session.Values)
	if err == nil {
		session.GetCache("cookie").(*cookie.Cookies).Set(session.Name(), val, opts)
	}
	return
}
