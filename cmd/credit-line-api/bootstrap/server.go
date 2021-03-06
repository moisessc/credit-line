package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"credit-line/pkg/env"
)

// server represents the server of the application
type server struct {
	router http.Handler
	Srv    *env.Server
}

// newServer create a new pointer of the server struct
func newServer(router http.Handler, srv *env.Server) *server {
	return &server{
		router: router,
		Srv:    srv,
	}
}

// up starts the HTTP server
func (s *server) up() error {

	srvPort := s.Srv.Port
	srvShutdownTimeOut := s.Srv.ShutdownTimeOut

	srvErr := make(chan error, 1)
	srvShutdown := make(chan os.Signal, 1)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", srvPort),
		Handler: s.router,
	}

	go func() {
		log.Printf("Server online on port: %v", srvPort)
		srvErr <- srv.ListenAndServe()
	}()

	signal.Notify(srvShutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		return fmt.Errorf("server error: %w", err)
	case shutdownSignal := <-srvShutdown:
		log.Println("starting shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(srvShutdownTimeOut))
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("gracefully shutdown failed")
			err = srv.Close()
		}

		switch {
		case shutdownSignal == syscall.SIGINT:
			log.Println("the stop signal caused shutdown")
		case err != nil:
			log.Println("could not stop server gracefully: %w", err)
		}

		log.Println("server shutdown...")
	}
	return nil
}
