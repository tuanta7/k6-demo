package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Run() error
	Shutdown(context.Context) error
}

func RunServers(servers []Server, gracePeriod ...time.Duration) error {
	errCh := make(chan error)
	for _, server := range servers {
		go func(s Server) {
			if err := s.Run(); err != nil {
				err = fmt.Errorf("error starting REST server: %w", err)
				errCh <- err
			}
		}(server)
	}

	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Println("Shutting down due to server error:", err)
		return shutdownServer(servers, gracePeriod...)
	case <-notifyCh:
		log.Println("Shutting down gracefully...")
		_ = shutdownServer(servers, gracePeriod...)
		return nil
	}
}

func shutdownServer(servers []Server, gracePeriod ...time.Duration) (err error) {
	timeout := 30 * time.Second
	if len(gracePeriod) > 0 {
		timeout = gracePeriod[0]
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, server := range servers {
		if se := server.Shutdown(shutdownCtx); se != nil {
			log.Println("Error during server shutdown:", se)
			err = errors.Join(err, se)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
