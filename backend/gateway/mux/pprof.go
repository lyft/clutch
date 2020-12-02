package mux

import (
	"net/http"
	"net/http/pprof"
)

func pprofHandlers(httpMux *http.ServeMux) {
	httpMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	httpMux.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	httpMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	httpMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	httpMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	httpMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	httpMux.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	httpMux.Handle("/debug/heap", pprof.Handler("heap"))
	httpMux.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	httpMux.Handle("/debug/block", pprof.Handler("block"))
	httpMux.Handle("/debug/mutex", pprof.Handler("mutex"))
}
