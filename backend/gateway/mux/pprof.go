package mux

import (
	"net/http"
	"net/http/pprof"
)

func pprofHandlers(httpMux *http.ServeMux) {
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
}
