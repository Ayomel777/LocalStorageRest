package handler

import (
	"Lab3_KSIS/internal/storage"
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer rc.Close()
		io.Copy(w, rc)

	case http.MethodPut:
		err := s.storage.Write(path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case http.MethodPost:
		err := s.storage.Append(path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		err := s.storage.Delete(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case "COPY":
		dst := strings.TrimPrefix(r.Header.Get("Destination"), "/files/")
		err := s.storage.Copy(path, dst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case "MOVE":
		dst := strings.TrimPrefix(r.Header.Get("Destination"), "/files/")
		err := s.storage.Move(path, dst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
