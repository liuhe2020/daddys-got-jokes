package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ServeFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type Server struct {
	addr string
	db   DB
}

func NewServer(addr string, db DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()
	// api
	r.Get("/joke", makeHTTPHandleFunc(s.handleJokeRandom))
	r.Get("/joke/{id}", makeHTTPHandleFunc(s.handleJokesById))
	r.Get("/jokes", makeHTTPHandleFunc(s.handleJokes))
	// static
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	server := http.Server{
		Addr:    s.addr,
		Handler: r,
	}
	log.Printf("Daddy's Got Jokes is running on port %s", s.addr)

	return server.ListenAndServe()
}

func (s *Server) handleJokes(w http.ResponseWriter, r *http.Request) error {
	page, err := getPage(r)
	if err != nil {
		return err
	}
	jokes, err := s.db.GetJokes(page)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, jokes)
}

func (s *Server) handleJokesById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}
	joke, err := s.db.GetJokeById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, joke)
}

func (s *Server) handleJokeRandom(w http.ResponseWriter, r *http.Request) error {
	joke, err := s.db.GetJokeRandom()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, joke)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f ServeFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: err.Error()})
		}
	}
}

func getPage(r *http.Request) (int, error) {
	pageDefault := 1
	pageStr := r.URL.Query().Get("page")
	// Validate params
	if len(r.URL.Query()) > 1 || (len(r.URL.Query()) == 1 && pageStr == "") {
		return 0, fmt.Errorf("invalid query param: only 'page' param is allowed")
	}
	// If 'page' parameter is not provided, return the default page
	if pageStr == "" {
		return pageDefault, nil
	}
	pageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid page number: %s", pageStr)
	}
	if pageNum < 1 {
		return 0, fmt.Errorf("page number must be greater than 0")
	}
	return pageNum, nil
}
