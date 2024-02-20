package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (s *Service) handleReq(ctx context.Context, request TaskCreateReq) error {
	// TODO: add request pool
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, nil)
	if err != nil {
		return fmt.Errorf("create new request: %w", errors.Join(ErrBadRequest, err))
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("do a http request: %w", err)
	}

	s.log.Debug("request roundtrip",
		slog.String("m", request.Method),
		slog.String("u", request.Url),
		slog.Int("h", len(request.Headers)),
		slog.Int("status", resp.StatusCode))

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
