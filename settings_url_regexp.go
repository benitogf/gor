package main

import (
	"net/http"
	"regexp"
)

// HTTPUrlRegexp ...
type HTTPUrlRegexp struct {
	regexp *regexp.Regexp
}

// String ...
func (r *HTTPUrlRegexp) String() string {
	if r.regexp == nil {
		return ""
	}
	return r.regexp.String()
}

// Set ...
func (r *HTTPUrlRegexp) Set(value string) error {
	regexp, err := regexp.Compile(value)
	r.regexp = regexp
	return err
}

// Good ...
func (r *HTTPUrlRegexp) Good(req *http.Request) bool {
	if r.regexp == nil {
		return true
	}
	return r.regexp.MatchString(req.Host + req.URL.String())
}
