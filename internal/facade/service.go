package facade

import (
	"context"
	"errors"
	"fmt"
	"github.com/kazhuravlev/sample-server/internal/api"
	"github.com/kazhuravlev/sample-server/pkg/bpool"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	log        *slog.Logger
	api        *api.Service
	port       int
	wg         *sync.WaitGroup
	bufferPool *bpool.Pool
}

func New(log *slog.Logger, apiInst *api.Service, port int) (*Service, error) {
	return &Service{
		log:        log,
		api:        apiInst,
		port:       port,
		wg:         new(sync.WaitGroup),
		bufferPool: bpool.New(),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	address := ":" + strconv.Itoa(s.port)

	handler := http.NewServeMux()
	handler.HandleFunc("/task", s.handleTaskCreate)

	httpServer := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadTimeout:       100 * time.Millisecond,
		ReadHeaderTimeout: 5 * time.Millisecond,
		WriteTimeout:      100 * time.Millisecond,
		IdleTimeout:       3 * time.Second,
		MaxHeaderBytes:    1024, // 1KB
		BaseContext:       func(lis net.Listener) context.Context { return ctx },
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		<-ctx.Done()

		// We should give more time when parent context was closed.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			s.log.Error("close http server", slog.String("error", err.Error()))
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil {
		// This is not an error in our case.
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run http server listener: %w", err)
		}
	}

	return nil
}

func (s *Service) Wait() {
	s.wg.Wait()
}
