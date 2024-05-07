package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	db         DB
}

type serveFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func NewServer(listenAddr string, db DB) *Server {
	return &Server{
		listenAddr: listenAddr,
		db:         db,
	}
}

func (s *Server) Run() {
	// rate limiter - 100/day & 10/minute
	limiter := tollbooth.NewLimiter(0.0694444, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetBurst(10)

	router := mux.NewRouter()
	// api
	router.Handle("/jokes", tollbooth.LimitHandler(limiter, makeHTTPHandleFunc(s.handleJokes)))
	router.Handle("/joke/{id}", tollbooth.LimitHandler(limiter, makeHTTPHandleFunc(s.handleJokesById)))
	router.Handle("/joke", tollbooth.LimitHandler(limiter, makeHTTPHandleFunc(s.handleJokeRandom)))
	router.Handle("/year", makeHTTPHandleFunc(s.handleYear))
	// static
	fs := http.FileServer(http.Dir("public"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	log.Println("JSON API server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func (s *Server) handleJokes(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
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
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleJokesById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getId(r)
		if err != nil {
			return err
		}
		joke, err := s.db.GetJokeById(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, joke)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleJokeRandom(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		joke, err := s.db.GetJokeRandom()
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, joke)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleYear(w http.ResponseWriter, r *http.Request) error {
	currentYear := strconv.Itoa(time.Now().Year())
	return WriteJSON(w, http.StatusOK, currentYear)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f serveFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
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
