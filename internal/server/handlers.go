package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
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
