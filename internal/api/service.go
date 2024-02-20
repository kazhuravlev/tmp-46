package api

import (
	"context"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	log            *slog.Logger
	validate       *validator.Validate
	http           *http.Client
	maxConcurrency int
	requests       chan TaskCreateReq
	semaphore      *semaphore.Weighted
}

func New(logger *slog.Logger) (*Service, error) {
	const maxConcurrency = 3

	return &Service{
		log:      logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
		http: &http.Client{
			Transport: &http.Transport{
				// TODO: configuration parameters. Tune it depends on workload.
				MaxIdleConns:        1024,
				MaxIdleConnsPerHost: 1024,
				MaxConnsPerHost:     1024,
				IdleConnTimeout:     3 * time.Second,
			},
			CheckRedirect: nil,
			Jar:           nil,
			// TODO: configuration parameter.
			Timeout: 10 * time.Second,
		},
		maxConcurrency: maxConcurrency,
		requests:       make(chan TaskCreateReq),
		semaphore:      semaphore.NewWeighted(int64(maxConcurrency)),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case req, ok := <-s.requests:
				if !ok {
					break Loop
				}

				if err := s.semaphore.Acquire(ctx, 1); err != nil {
					log.Printf("Failed to acquire semaphore: %v", err)
					break
				}

				go func(req TaskCreateReq) {
					defer s.semaphore.Release(1)

					if err := s.handleReq(ctx, req); err != nil {
						s.log.Error("handle request", slog.String("error", err.Error()))
					}
				}(req)
			}
		}
	}()

	return nil
}

func (s *Service) Wait() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.semaphore.Acquire(ctx, int64(s.maxConcurrency)); err != nil {
		s.log.Error("acquire worker semaphore", slog.String("error", err.Error()))
	}
}
