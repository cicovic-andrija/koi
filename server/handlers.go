package server

import "net/http"

func multiplexer() http.Handler {
	mux := http.NewServeMux()

	return mux
}
