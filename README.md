# cookie-session

Use cookie as session, base on [secure cookie](https://github.com/go-http-utils/cookie) to encrypt cookie, you can also use the another library instead of it.

[![Build Status](https://travis-ci.org/go-http-utils/cookie-session.svg?branch=master)](https://travis-ci.org/go-http-utils/cookie-session)
[![Coverage Status](http://img.shields.io/coveralls/go-http-utils/cookie-session.svg?style=flat-square)](https://coveralls.io/r/go-http-utils/cookie-session)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-http-utils/cookie-session/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-http-utils/cookie-session)

## Features

* Simple API: use it as an easy way to set signed cookies.
* Built-in backends to store sessions in cookies.
* Mechanism to rotate authentication by some custom keys.
* Multiple sessions per request, even using different backends.
* Interfaces and infrastructure for custom session backends: sessions from
  different stores can be retrieved and batch-saved using a common API.
* User can customize own session with different field that don't require type assertion and cast

## Installation

```go
go get github.com/go-http-utils/cookie-session
```

## Examples

```go
go run example/main.go
```

## Usage

```go
// Session is custom by user's business
type Session struct {
  *sessions.Meta `json:"-"`
  UserID         string `json:"userId"`
  Name           string `json:"name"`
  Authed         int64  `json:"authed"`
}

func (s *Session) Save() error {
  return s.GetStore().Save(s)
}

func (s *Session) Destroy() error {
  return s.GetStore().Destroy(s)
}

func main() {
  SessionName := "Sess"
  SessionKeys := []string{"keyxxx"}

  store := sessions.New()

  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      session := &Session{Meta: &sessions.Meta{}}
    store.Load(SessionName, session, cookie.New(w, r, SessionKeys...))
    if session.UserID == "" {
      session.UserID = "x"
      session.Name = "y"
      session.Authed = 1
      session.Save()
    }
  })
```

## Other Store Implementations

* https://github.com/mushroomsir/session-redis -Redis

## Other Application Implementations

* https://github.com/teambition/gear-session