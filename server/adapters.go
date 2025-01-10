package server

import "net/http"

// Adapter is an HTTP(S) handler that invokes another HTTP(S) handler.
type Adapter func(h http.Handler) http.Handler

// Adapt returns an HTTP(S) handler enhanced by a number of adapters.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// StripPrefix returns an adapter that calls http.StripPrefix
// to remove the given prefix from the request's URL path and invoke
// the handler h.
func StripPrefix(prefix string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	}
}
