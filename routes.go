package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
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

	// r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	w.Write([]byte("method is not valid"))
	// })

	// api
	r.Group(func(r chi.Router) {
		r.Use(httprate.Limit(
			100,
			24*time.Hour,
			httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "maximum per-minute requests reached, try again later", http.StatusTooManyRequests)
			}),
		))

		r.Route("/joke", func(r chi.Router) {
			r.Get("/", s.handleJokeRandom)
			r.Get("/{id}", s.handleJokesById)
		})

		r.Get("/jokes", s.handleJokes)
	})

	// static
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	server := http.Server{
		Addr:    s.addr,
		Handler: r,
	}
	log.Printf("Daddy's Got Jokes is running on port %s", s.addr)

	return server.ListenAndServe()
}

func (s *Server) handleJokes(w http.ResponseWriter, r *http.Request) {
	page, err := getPage(r)
	if err != nil {
		s.handleError(w, http.StatusBadRequest, err)
		return
	}
	jokes, err := s.db.GetJokes(page)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSON(w, http.StatusOK, jokes)
}

func (s *Server) handleJokesById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("failed to parse string id to an integer from url param")
		return
	}
	joke, err := s.db.GetJokeById(id)
	if err != nil {
		s.handleError(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, joke)
}

func (s *Server) handleJokeRandom(w http.ResponseWriter, r *http.Request) {
	joke, err := s.db.GetJokeRandom()
	if err != nil {
		s.handleError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSON(w, http.StatusOK, joke)
}

func (s *Server) handleError(w http.ResponseWriter, status int, err error) {
	log.Printf("error: %v", err)
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error encoding response: %v", err)
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
