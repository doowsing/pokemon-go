package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	DefaultKey  = "github.com/gin-contrib/sessions"
	errorFormat = "[sessions] ERROR! %s\n"
)

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) (interface{}, bool)
	// Get returns the session value associated to the given key.
	MustGet(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the session.
	Clear()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error

	Session() *sessions.Session

	SessionId() string
}

func Sessions(name string, store sessions.Store, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

func SessionsMany(names []string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		_sessions := make(map[string]Session, len(names))
		for _, name := range names {
			_sessions[name] = &session{name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(DefaultKey, _sessions)
		defer context.Clear(c.Request)
		c.Next()
	}
}

type session struct {
	name    string
	request *http.Request
	store   sessions.Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) MustGet(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Get(key interface{}) (interface{}, bool) {
	result, ok := s.Session().Values[key]
	return result, ok
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options Options) {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (s *session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *session) Written() bool {
	return s.written
}

func (s *session) SessionId() string {
	return s.Session().ID
}

// shortcut to get session
func GetSession(c *gin.Context) Session {
	session, ok := c.Get(DefaultKey)
	if !ok {
		return nil
	}
	trueSession, ok := session.(Session)
	if !ok {
		return nil
	}
	return trueSession
}

// shortcut to get session with given name
func DefaultMany(c *gin.Context, name string) Session {
	return c.MustGet(DefaultKey).(map[string]Session)[name]
}
