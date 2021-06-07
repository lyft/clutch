package mux

import (
	"net/http"
	"reflect"
	"strings"
)

// All modern browsers send `text/html` as the first MIME type in the Accept header for navigation requests.
// This is naive and doesn't do any actual content negotiation, just checks if the client accepts text/html which
// is good enough for our use case of determining whether the request originates from the SPA's XHR network client or
// browser navigation.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation/List_of_default_Accept_values
func isBrowser(h http.Header) bool {
	directives := strings.Split(h.Get("Accept"), ",")
	for _, d := range directives {
		mt := strings.SplitN(strings.TrimSpace(d), ";", 1)
		if len(mt) > 0 && mt[0] == "text/html" {
			return true
		}
	}
	return false
}

// This is a hack, but it's the best way to get the Accept header from the response forwarder without significant plumbing
// complexity and overhead via middleware. Do not invoke this in the normal request flow.
func requestHeadersFromResponseWriter(w http.ResponseWriter) http.Header {
	req := reflect.ValueOf(w).Elem().FieldByName("req")
	if !req.IsValid() {
		return nil
	}
	h := req.Elem().FieldByName("Header")
	if !h.IsValid() {
		return nil
	}

	ret := make(http.Header, h.Len())
	iter := h.MapRange()
	for iter.Next() {
		k := iter.Key().String()
		var v string
		if iter.Value().Len() > 0 {
			v = iter.Value().Index(0).String()
		}
		ret[k] = []string{v}
	}
	return ret
}

// GetCookieValue is the easiest way to parse a cookie string in a non-HTTP request context.
func GetCookieValue(headerValues []string, key string) (string, error) {
	if key == "" {
		return "", http.ErrNoCookie
	}

	request := http.Request{Header: http.Header{"Cookie": headerValues}}
	c, err := request.Cookie(key)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
