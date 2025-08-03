package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"

	"gosuda.org/boilerplate/api"
	"gosuda.org/boilerplate/pkg/blog"
	"gosuda.org/boilerplate/pkg/config"
	"gosuda.org/boilerplate/pkg/debug"
	"gosuda.org/boilerplate/pkg/logging"
	mw "gosuda.org/boilerplate/pkg/middleware"
	"gosuda.org/boilerplate/pkg/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	buf := &logging.Buffer{}
	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, buf), nil))

	r := chi.NewRouter()
	r.Use(chimid.RequestID)
	r.Use(mw.Logger(logger))
	r.Use(mw.Recover)

	st := store.NewMemory()
	srv := blog.New(st)
	handler := api.NewStrictHandler(srv, nil)
	api.HandlerFromMux(handler, r)
	r.Mount("/debug", debug.Routes(buf.Lines))

	server := &http.Server{Addr: cfg.Server.Addr, Handler: r}
	go func() {
		logger.Info("starting server", "addr", cfg.Server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	_ = st.Close()
	logger.Info("server shutdown")
}
