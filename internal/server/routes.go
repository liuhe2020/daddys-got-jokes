package server

import (
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/liuhe2020/daddys-got-jokes/cmd/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// api
	r.Group(func(r chi.Router) {
		r.Use(httprate.Limit(
			100,
			24*time.Hour,
			httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "maximum per-minute requests reached, try again later", http.StatusTooManyRequests)
			}),
		))

		// r.Route will check for method and respond with appropriate error msg
		r.Route("/joke", func(r chi.Router) {
			r.Get("/", s.handleJokeRandom)
			r.Get("/{id}", s.handleJokesById)
		})

		r.Route("/jokes", func(r chi.Router) {
			r.Get("/", s.handleJokes)
		})
	})

	r.Get("/", templ.Handler(web.Base()).ServeHTTP)
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/health", s.healthHandler)
	// static
	// r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	log.Printf("Daddy's Got Jokes is running on port %d", s.port)
	return r
}
