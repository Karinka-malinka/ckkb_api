package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	HTTPServer      http.Server
	ShutdownTimeout time.Duration
}

func NewHTTPServer(addr string, shutdownTimeout time.Duration, h http.Handler) *Server {

	server := &Server{
		ShutdownTimeout: shutdownTimeout,
	}

	server.HTTPServer = http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server
}

func (s *Server) Start(ctx context.Context) {

	errs, _ := errgroup.WithContext(ctx)

	errs.Go(func() error {
		logrus.Info("HTTPS server starting on ", s.HTTPServer.Addr)

		err := s.HTTPServer.ListenAndServe()

		// игнорируем нормальное завершение
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	})

	err := errs.Wait()
	if err != nil {
		logrus.Errorf("message from server: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) {

	ctx, cancelShutdown := context.WithTimeout(ctx, s.ShutdownTimeout*time.Second)
	defer cancelShutdown()

	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("server shutdown with error: %v", err)
	}

	logrus.Infof("Server is graceful shutdown...")
}
