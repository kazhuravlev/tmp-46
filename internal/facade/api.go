package facade

import (
	"encoding/json"
	"github.com/kazhuravlev/sample-server/internal/api"
	"io"
	"log/slog"
	"net/http"
)

type TaskCreateReq struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type TaskCreateResp struct {
	ID string `json:"id"`
}

func (s *Service) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.log.Error("bad request method", slog.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req TaskCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Error("bad request", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiReq := api.TaskCreateReq{
		Method:  req.Method,
		Url:     req.Url,
		Headers: req.Headers,
	}
	id, err := s.api.TaskCreate(r.Context(), apiReq)
	if err != nil {
		s.log.Error("create task", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := TaskCreateResp{
		ID: id,
	}
	buf := s.bufferPool.Get()
	defer s.bufferPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		s.log.Error("enode response", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, buf); err != nil {
		s.log.Error("write response", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
