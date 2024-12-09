package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pdkonovalov/auth-server/pkg/config"
	"github.com/pdkonovalov/auth-server/pkg/email"
	v1 "github.com/pdkonovalov/auth-server/pkg/http/api/v1"
	"github.com/pdkonovalov/auth-server/pkg/jwt"
	"github.com/pdkonovalov/auth-server/pkg/storage"
)

type Server struct {
	httpServer *http.Server
}

func MakeServer(config *config.Config, storage storage.Storage, email *email.Email, jwt *jwt.JwtGenerator) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Route("/api/v1/jwt", func(r chi.Router) {
		r.Get("/new", v1.HandleNewJwt(storage, jwt))
		r.Post("/refresh", v1.HandleRefreshJwt(storage, email, jwt))
	})
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: r,
	}
	return &Server{
		httpServer: httpServer,
	}
}

func (srv *Server) Start() error {
	go func() {
		err := srv.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	return nil
}

func (srv *Server) Shutdown() error {
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()
	err := srv.httpServer.Shutdown(shutdownCtx)
	return err
}
