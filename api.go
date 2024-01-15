package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	db         DB
}

func NewAPIServer(listenAddr string, db DB) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		db:         db,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/jokes", makeHTTPHandleFunc(s.handleJokes))
	router.HandleFunc("/joke/{id}", makeHTTPHandleFunc(s.handleJokesById))
	// router.HandleFunc("/joke/random", makeHTTPHandleFunc(s.handleJokesRandom))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleJokes(w http.ResponseWriter, r *http.Request) error {
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

func (s *APIServer) handleJokesById(w http.ResponseWriter, r *http.Request) error {
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

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
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
