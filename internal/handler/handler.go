package handler

import (
	"Lab3_KSIS/internal/storage"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Server struct {
	storage *storage.Storage
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{storage: storage}
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/files/")

	switch r.Method {
	case http.MethodGet:
		rc, err := s.storage.Read(path)
		if err != nil {
			writeErr(w, err)
			return
		}
		defer rc.Close()
		io.Copy(w, rc)

	case http.MethodPut:
		writeErr(w, s.storage.Write(path, r.Body))

	case http.MethodPost:
		writeErr(w, s.storage.Append(path, r.Body))

	case http.MethodDelete:
		writeErr(w, s.storage.Delete(path))

	case "COPY":
		dst := strings.TrimPrefix(r.Header.Get("Destination"), "/files/")
		writeErr(w, s.storage.Copy(path, dst))

	case "MOVE":
		dst := strings.TrimPrefix(r.Header.Get("Destination"), "/files/")
		writeErr(w, s.storage.Move(path, dst))

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeErr(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, storage.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, storage.ErrBadPath):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
