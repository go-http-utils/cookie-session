# cookie-session
Use cookie as session, base on [secure cookie](https://github.com/go-http-utils/cookie) to encrypt cookie, you can also use the another library instead of it.

[![Build Status](https://travis-ci.org/go-http-utils/cookie-session.svg?branch=master)](https://travis-ci.org/go-http-utils/cookie-session)
[![Coverage Status](http://img.shields.io/coveralls/go-http-utils/cookie-session.svg?style=flat-square)](https://coveralls.io/r/go-http-utils/cookie-session)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-http-utils/cookie-session/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-http-utils/cookie-session)

##Features
* Simple API: use it as an easy way to set signed cookies.
* Built-in backends to store sessions in cookies.
* Mechanism to rotate authentication by some custom keys.
* Multiple sessions per request, even using different backends.
* Interfaces and infrastructure for custom session backends: sessions from
  different stores can be retrieved and batch-saved using a common API.

##Examples
```go
go run cookiesession/main.go
```
##Usage
```go
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := cookiesession.New(w, r)
		session, _ := store.Get("login")
		session.Values["name"] = "mushroom"
		session.Values[66] = 99  //set session
		println(session.Values["name"].(string)) //get session
		session.Save()
	})
```
