package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lukejoshuapark/app"
)

type Service struct {
	port    uint16
	handler http.Handler

	configureServer func(*http.Server)
	shutdownTimeout time.Duration
	certFile        string
	keyFile         string
}

var _ app.Service = (*Service)(nil)

func NewService(port uint16, handler http.Handler) *Service {
	return &Service{
		port:    port,
		handler: handler,

		shutdownTimeout: 10 * time.Second,
	}
}

func (s *Service) WithServerConfiguration(fn func(*http.Server)) *Service {
	s.configureServer = fn
	return s
}

func (s *Service) WithShutdownTimeout(timeout time.Duration) *Service {
	s.shutdownTimeout = timeout
	return s
}

func (s *Service) WithTLS(certFile, keyFile string) *Service {
	s.certFile = certFile
	s.keyFile = keyFile
	return s
}

func (s *Service) Run(ctx context.Context) error {
	svr := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           s.handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if s.configureServer != nil {
		s.configureServer(svr)
	}

	errc := make(chan error, 1)
	go func() {
		if s.certFile != "" && s.keyFile != "" {
			errc <- svr.ListenAndServeTLS(s.certFile, s.keyFile)
		} else {
			errc <- svr.ListenAndServe()
		}
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := svr.Shutdown(ctx); err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
