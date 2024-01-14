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
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/joke", makeHTTPHandleFunc(s.handleJoke))
	router.HandleFunc("/joke/{id}", makeHTTPHandleFunc(s.handleJokeById))
	// router.HandleFunc("/joke/random", makeHTTPHandleFunc(s.handleJokeRandom))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleJoke(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		page, err := getPage(r)
		if err != nil {
			return err
		}
		jokes, err := s.store.GetJokes(page)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, jokes)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleJokeById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getId(r)
		if err != nil {
			return err
		}

		joke, err := s.store.GetJokeById(id)
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
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		return 1, nil // default page
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid page given %s", pageStr)
	}
	return page, nil
}
