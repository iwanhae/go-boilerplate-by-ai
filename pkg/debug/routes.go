package debug

import (
	"expvar"
	"io"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

// Routes returns a router with debug endpoints.
func Routes(logs func() []string) http.Handler {
	r := chi.NewRouter()
	r.Get("/metrics", expvar.Handler().ServeHTTP)
	r.Get("/logs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		for _, l := range logs() {
			io.WriteString(w, l+"\n")
		}
	})

	p := chi.NewRouter()
	p.Get("/", pprof.Index)
	p.Get("/cmdline", pprof.Cmdline)
	p.Get("/profile", pprof.Profile)
	p.Get("/symbol", pprof.Symbol)
	p.Post("/symbol", pprof.Symbol)
	p.Get("/trace", pprof.Trace)
	p.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
	p.Get("/block", pprof.Handler("block").ServeHTTP)
	p.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	p.Get("/heap", pprof.Handler("heap").ServeHTTP)
	p.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
	p.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	r.Mount("/pprof", p)
	return r
}
