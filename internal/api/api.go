package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"net/http"
	"strings"
	"time"
)

// TODO: Depends on purposes of this server we should fill this map more accurate or change to blacklist style
// 	or even remove it.

var allowedHeaders = map[string]struct{}{
	"authorization": {},
	"cookie":        {},
}

type TaskCreateReq struct {
	Method  string `validate:"required"`
	Url     string `validate:"required,url"`
	Headers map[string]string
}

// TaskCreate will create and run new request task.
// NOTE: This function do not accept the body. See requirements.
func (s *Service) TaskCreate(ctx context.Context, req TaskCreateReq) (string, error) {
	if err := s.validate.Struct(req); err != nil {
		return "", fmt.Errorf("bad request: %w", errors.Join(ErrBadRequest, err))
	}

	// NOTE: we can write a custom validation function for this field.
	switch req.Method {
	default:
		return "", fmt.Errorf("unknown method: %w", ErrBadRequest)
	case http.MethodPost, http.MethodGet:
	}

	for headerName := range req.Headers {
		if _, ok := allowedHeaders[strings.ToLower(headerName)]; !ok {
			return "", fmt.Errorf("not allowed headers: %w", ErrBadRequest)
		}
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("context cancelled: %w", ctx.Err())
	case <-time.After(20 * time.Millisecond): // TODO: tune this.
		return "", fmt.Errorf("have no capacity to add register new task: %w", ErrInternal)
	case s.requests <- req:
	}

	return ulid.Make().String(), nil
}
