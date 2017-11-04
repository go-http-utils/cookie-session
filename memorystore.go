package sessions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sync"
	"time"

	"github.com/go-http-utils/cookie"
)

// NewMemoryStore returns an MemoryStore instance
func NewMemoryStore(options ...*Options) (store *MemoryStore) {
	opts := &cookie.Options{
		Path:     "/",
		HTTPOnly: true,
		Signed:   false, // not necessary
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
	store = &MemoryStore{
		opts:   opts,
		ticker: time.NewTicker(time.Second),
		store:  make(map[string]*sessionValue),
		done:   make(chan bool, 1),
	}

	go store.cleanCache()
	return
}

type sessionValue struct {
	expired time.Time
	session string
}

// MemoryStore using memory to store sessions base on secure cookies.
type MemoryStore struct {
	opts   *cookie.Options
	store  map[string]*sessionValue
	ticker *time.Ticker
	lock   sync.Mutex
	done   chan bool
}

// Load a session by name and any kind of stores
func (m *MemoryStore) Load(name string, session Sessions, cookie *cookie.Cookies) error {
	sid, err := cookie.Get(name, m.opts.Signed)
	var result string
	if sid != "" {
		m.lock.Lock()
		if val, ok := m.store[sid]; ok {
			result = val.session
		}
		m.lock.Unlock()
	}
	if result != "" {
		err = Decode(result, &session)
	}
	session.Init(name, sid, cookie, m, result)
	return err
}

// Save session to Response's cookie
func (m *MemoryStore) Save(session Sessions) (err error) {
	val, err := Encode(session)
	if err != nil || !session.IsChanged(val) {
		return
	}
	sid := session.GetSID()
	if sid == "" {
		sid = NewSID(val)
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.store[sid] = &sessionValue{
		session: val,
		expired: time.Now().Add(time.Duration(m.opts.MaxAge) * time.Second),
	}
	session.GetCookie().Set(session.GetName(), sid, m.opts)
	return
}

// Destroy destroy the session
func (m *MemoryStore) Destroy(session Sessions) (err error) {
	sid := session.GetSID()
	if sid != "" {
		m.lock.Lock()
		defer m.lock.Unlock()
		delete(m.store, sid)
	}
	session.GetCookie().Remove(session.GetName(), m.opts)
	return
}

// Len ...
func (m *MemoryStore) Len() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.store)
}

// Close goroutine cleanCache thread
func (m *MemoryStore) Close() {
	close(m.done)
}

func (m *MemoryStore) cleanCache() {
	defer m.ticker.Stop()
	for {
		select {
		case <-m.ticker.C:
			m.clean()
		case <-m.done:
			return
		}
	}
}

func (m *MemoryStore) clean() {
	m.lock.Lock()
	defer m.lock.Unlock()
	start := time.Now()
	expireTime := start.Add(time.Millisecond * 100)
	frequency := 24
	var expired int
	for {
	label:
		for i := 0; i < frequency; i++ {
			for key, value := range m.store {
				if value.expired.Before(start) {
					delete(m.store, key)
					expired++
				}
				break
			}
		}
		if expireTime.Before(time.Now()) {
			return
		}
		if expired > frequency/4 {
			expired = 0
			goto label
		}
		return
	}
}

// NewSID generates a random SID
func NewSID(val string) string {
	h := sha256.New()
	h.Write([]byte(val))
	io.CopyN(h, rand.Reader, 8)
	return hex.EncodeToString(h.Sum(nil))
}
