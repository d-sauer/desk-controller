package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d-sauer/exploring-go/desk-controller/internal/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func main() {
	var port string

	flag.StringVar(&port, "port", "9696", "HTTP Server Address")
	flag.Parse()

	errC, err := run(port)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(port string) (chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("can't initialize logger: %w", err)
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}

	srv, err := newServer(serverConfig{
		Address:     ":" + port,
		Logger:      logger,
		Middlewares: []func(next http.Handler) http.Handler{logging},
	})

	if err != nil {
		return nil, fmt.Errorf("can't initialize server (port: %s): %w", port, err)
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			_ = logger.Sync()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil { //nolint: contextcheck
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("port", port))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, err
}

type serverConfig struct {
	Address     string
	Middlewares []func(next http.Handler) http.Handler
	Logger      *zap.Logger
}

func newServer(config serverConfig) (*http.Server, error) {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))

	// A good base middleware stack
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	// Load middlewares
	for _, mw := range config.Middlewares {
		router.Use(mw)
	}

	router.Mount("/", rest.DeskControllerRouter())
	router.Mount("/app", appRouter())

	return &http.Server{
		Handler:           router,
		Addr:              config.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
