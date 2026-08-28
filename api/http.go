package api

import (
	"drupal-scheduler/model"
	"drupal-scheduler/service"
	"encoding/json"
	"net/http"
)

type Server struct{ svc *service.Service }

func New(s *service.Service) *Server { return &Server{svc: s} }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	return mux
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rec model.Record
		if json.NewDecoder(r.Body).Decode(&rec) != nil {
			http.Error(w, "bad json", 400)
			return
		}
		if e := s.svc.Register(r.Context(), rec); e != nil {
			http.Error(w, e.Error(), 422)
			return
		}
		w.WriteHeader(201)
		return
	}
	rows, e := s.svc.Query(r.URL.Query().Get("status"), r.URL.Query().Get("q"))
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(rows)
}
